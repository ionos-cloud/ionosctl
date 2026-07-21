package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/globalwait"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunServerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	return nil
}

func PreRunServerCreate(c *core.PreCommandConfig) error {
	serverType := c.Flags().String(constants.FlagType)
	requiredFlags, err := getRequiredFlagsByServerType(serverType)
	if err != nil {
		return err
	}

	// Confidential VMs are ENTERPRISE-only, and their cores + CPU family are derived from the boot
	// image's launch-config.json — the API rejects them on the request, so they must not be set here.
	if c.Flags().Bool(constants.FlagConfidential) {
		if serverType != serverEnterpriseType {
			return fmt.Errorf("--%s requires --%s %s (Confidential VMs are ENTERPRISE-only)",
				constants.FlagConfidential, constants.FlagType, serverEnterpriseType)
		}
		changed := c.Command.Command.Flags().Changed
		if changed(constants.FlagCores) || changed(constants.FlagCpuFamily) {
			return fmt.Errorf("--%s: do not set --%s or --%s; both are derived from the confidential image's launch-config.json",
				constants.FlagConfidential, constants.FlagCores, constants.FlagCpuFamily)
		}
		if !changed(cloudapiv6.ArgImageId) {
			return fmt.Errorf("--%s requires --%s: a Confidential VM must boot from a confidential image "+
				"(find one with: ionosctl image list -F public=false,requiredFeatures=SEV-SNP)",
				constants.FlagConfidential, cloudapiv6.ArgImageId)
		}
		// cores is image-derived here, so drop it from the ENTERPRISE required set.
		filtered := make([]string, 0, len(requiredFlags))
		for _, f := range requiredFlags {
			if f != constants.FlagCores {
				filtered = append(filtered, f)
			}
		}
		requiredFlags = filtered
	}

	if err = core.CheckRequiredFlags(c.Command, c.NS, requiredFlags...); err != nil {
		return fmt.Errorf("missing %s flags: %w", serverType, err)
	}

	// CUBE Attached Volume promotion logic (--promote-volume)
	// --promote-volume requires --wait and --type CUBE/GPU to be set
	if c.Flags().Bool(constants.FlagPromoteVolume) {
		if !viper.GetBool(constants.ArgWait) {
			return fmt.Errorf("--%s requires --%s to be set", constants.FlagPromoteVolume, constants.ArgWait)
		}
		serverType := c.Flags().String(constants.FlagType)
		if serverType != serverCubeType && serverType != serverGPUType {
			return fmt.Errorf("--%s can only be used with --%s %s or %s",
				constants.FlagPromoteVolume, constants.FlagType, serverCubeType, serverGPUType)
		}
	}

	// Check if image ID or alias is set
	if c.Flags().Changed(cloudapiv6.ArgImageId) || c.Flags().Changed(cloudapiv6.ArgImageAlias) {
		imageRequiredFlags := make([][]string, 0)

		if c.Flags().Changed(cloudapiv6.ArgImageAlias) {
			// Handle public image alias
			imageRequiredFlags = append(imageRequiredFlags,
				append(requiredFlags, cloudapiv6.ArgImageAlias, cloudapiv6.ArgPassword),
				append(requiredFlags, cloudapiv6.ArgImageAlias, cloudapiv6.ArgSshKeyPaths),
			)
		} else if c.Flags().Changed(cloudapiv6.ArgImageId) {
			// Check if the image ID corresponds to an image or snapshot
			img, _, imgErr := client.Must().CloudClient.ImagesApi.ImagesFindById(context.Background(),
				c.Flags().String(cloudapiv6.ArgImageId)).Execute()
			if imgErr != nil {
				// Try to fetch as a snapshot if image fetch fails
				_, _, snapshotErr := client.Must().CloudClient.SnapshotsApi.SnapshotsFindById(context.Background(),
					c.Flags().String(cloudapiv6.ArgImageId)).Execute()
				if snapshotErr != nil {
					return fmt.Errorf("failed getting image or snapshot %s: %w", c.Flags().String(cloudapiv6.ArgImageId), imgErr)
				}

				// If it's a snapshot, no additional checks are required
				return nil
			}

			// If it's an image, determine if it is public or private
			if img.Properties != nil && img.Properties.Public != nil && *img.Properties.Public {
				// For public images, require password or SSH key
				imageRequiredFlags = append(imageRequiredFlags,
					append(requiredFlags, cloudapiv6.ArgImageId, cloudapiv6.ArgPassword),
					append(requiredFlags, cloudapiv6.ArgImageId, cloudapiv6.ArgSshKeyPaths),
				)
			} else {
				// For private images, only the image ID is required
				imageRequiredFlags = append(imageRequiredFlags,
					append(requiredFlags, cloudapiv6.ArgImageId),
				)
			}
		}

		err = core.CheckRequiredFlagsSets(c.Command, c.NS, imageRequiredFlags...)
		if err != nil {
			return err
		}
	}

	return nil
}

func PreRunDcServerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId)
}

func PreRunDcServerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func RunServerListAll(c *core.CommandConfig) error {
	datacenters, _, err := c.CloudApiV6Services.DataCenters().List()
	if err != nil {
		return err
	}

	allDcs := helpers.GetDataCenters(datacenters)
	var allServers []ionoscloud.Servers
	totalTime := time.Duration(0)

	for _, dc := range allDcs {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("could not retrieve Datacenter Id")
		}

		servers, resp, err := c.CloudApiV6Services.Servers().List(*id)
		if err != nil {
			return err
		}

		allServers = append(allServers, servers.Servers)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer(AllServerCols).Prefix("*.items").Print(allServers)
}

func RunServerList(c *core.CommandConfig) error {
	if c.Flags().Bool(cloudapiv6.ArgAll) {
		return RunServerListAll(c)
	}

	servers, resp, err := c.CloudApiV6Services.Servers().List(c.Flags().String(cloudapiv6.ArgDataCenterId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(AllServerCols).Prefix("items").Print(servers.Servers)
}

func RunServerGet(c *core.CommandConfig) error {
	c.Verbose("Server with id: %v is getting... ", c.Flags().String(cloudapiv6.ArgServerId))

	svr, resp, err := c.CloudApiV6Services.Servers().Get(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(AllServerCols).Print(svr.Server)
}

func RunServerCreate(c *core.CommandConfig) error {
	input, err := getNewServer(c)
	if err != nil {
		return err
	}

	// If Server is of type CUBE, it will create an attached Volume
	if c.Flags().String(constants.FlagType) == serverCubeType {
		// Volume Properties
		volumeDAS, err := getNewDAS(c)
		if err != nil {
			return err
		}

		// Attach Storage
		input.SetEntities(ionoscloud.ServerEntities{
			Volumes: &ionoscloud.AttachedVolumes{
				Items: &[]ionoscloud.Volume{volumeDAS.Volume},
			},
		})
	}

	// If Server is of type GPU, it will create an attached Volume
	if c.Flags().String(constants.FlagType) == serverGPUType {
		// Volume Properties
		volumeGPU, err := getNewDAS(c)
		if err != nil {
			return err
		}

		// Attach Storage
		input.SetEntities(ionoscloud.ServerEntities{
			Volumes: &ionoscloud.AttachedVolumes{
				Items: &[]ionoscloud.Volume{volumeGPU.Volume},
			},
		})
	}

	// A Confidential VM must be created together with a boot volume built from the confidential
	// image, in the same request — the API derives cores + CPU family from that image.
	if c.Flags().Bool(constants.FlagConfidential) {
		volumeConf, err := getNewDAS(c)
		if err != nil {
			return err
		}

		input.SetEntities(ionoscloud.ServerEntities{
			Volumes: &ionoscloud.AttachedVolumes{
				Items: &[]ionoscloud.Volume{volumeConf.Volume},
			},
		})
	}

	svr, resp, err := c.CloudApiV6Services.Servers().Create(c.Flags().String(cloudapiv6.ArgDataCenterId), *input)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if c.Flags().Bool(constants.FlagPromoteVolume) {
		if err = promoteVolume(c, svr); err != nil {
			return err
		}
	}

	return c.Printer(AllServerCols).Print(svr.Server)
}

// promoteVolume waits for the server to reach AVAILABLE, then PATCHes it to set
// the first attached volume as the boot volume. Both waits (POST and PATCH) run
// under a single progress bar; MarkDone() prevents the post-command
// WaitAndRerender from showing a second bar.
func promoteVolume(c *core.CommandConfig, svr *resources.Server) error {
	id, ok := svr.GetIdOk()
	if !ok || id == nil {
		return errors.New("error getting new server id")
	}

	stderr := c.Command.Command.ErrOrStderr()
	cfg := client.Must().CloudClient.GetConfig()

	// One progress bar for the entire promote-volume operation.
	bar := pb.New(1)
	bar.SetWriter(stderr)
	bar.SetTemplateString(globalwait.ProgressTpl)
	bar.Start()

	// Phase 1: wait for server POST to reach AVAILABLE.
	// The capturing transport already stored the server href from POST.
	if err := globalwait.WaitForAvailable(io.Discard, cfg.Token, cfg.Username, cfg.Password); err != nil {
		bar.SetTemplateString(globalwait.ProgressTpl + " FAILED")
		bar.Finish()
		return err
	}
	globalwait.Reset()

	// Fetch fresh server data with entities populated.
	freshSvr, _, err := c.CloudApiV6Services.Servers().Get(
		c.Flags().String(cloudapiv6.ArgDataCenterId), *id)
	if err != nil {
		bar.SetTemplateString(globalwait.ProgressTpl + " FAILED")
		bar.Finish()
		return err
	}

	if freshSvr.Server.Entities == nil || freshSvr.Server.Entities.Volumes == nil ||
		freshSvr.Server.Entities.Volumes.Items == nil || len(*freshSvr.Server.Entities.Volumes.Items) == 0 {
		bar.SetTemplateString(globalwait.ProgressTpl + " FAILED")
		bar.Finish()
		return errors.New("no attached volumes found to promote to boot volume")
	}

	// PATCH: set boot volume to the first attached volume.
	attachedDas := (*freshSvr.Server.Entities.Volumes.Items)[0]
	bootVolume := ionoscloud.ResourceReference{Id: attachedDas.Id}
	updatedServer, _, err := client.Must().CloudClient.ServersApi.DatacentersServersPatch(
		context.Background(),
		c.Flags().String(cloudapiv6.ArgDataCenterId), *svr.Id).
		Server(ionoscloud.ServerProperties{BootVolume: &bootVolume}).Execute()
	if err != nil {
		bar.SetTemplateString(globalwait.ProgressTpl + " FAILED")
		bar.Finish()
		return fmt.Errorf("error promoting attached volume to boot volume: %w", err)
	}

	// Phase 2: wait for PATCH to reach AVAILABLE.
	if err := globalwait.WaitForAvailable(io.Discard, cfg.Token, cfg.Username, cfg.Password); err != nil {
		bar.SetTemplateString(globalwait.ProgressTpl + " FAILED")
		bar.Finish()
		return err
	}

	bar.SetTemplateString(globalwait.ProgressTpl + " DONE")
	bar.Finish()

	// All waiting done inline — skip post-command WaitAndRerender.
	globalwait.Reset()
	globalwait.MarkDone()

	// Re-fetch server with final AVAILABLE state so JSON output is fresh.
	freshSvr, _, err = c.CloudApiV6Services.Servers().Get(
		c.Flags().String(cloudapiv6.ArgDataCenterId), *id)
	if err != nil {
		// Non-fatal: fall back to PATCH response (may show BUSY state in JSON).
		svr.Server = updatedServer
		return nil
	}
	svr.Server = freshSvr.Server
	return nil
}

func RunServerUpdate(c *core.CommandConfig) error {
	input, err := getUpdateServerInfo(c)
	if err != nil {
		return err
	}

	c.Verbose("Updating Server with ID: %v in Datacenter with ID: %v",
		c.Flags().String(cloudapiv6.ArgServerId),
		c.Flags().String(cloudapiv6.ArgDataCenterId))

	svr, resp, err := c.CloudApiV6Services.Servers().Update(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
		*input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(AllServerCols).Print(svr.Server)
}

func RunServerDelete(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	serverId := c.Flags().String(cloudapiv6.ArgServerId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DeleteAllServers(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Server with id: %v from datacenter with id: %v... ", serverId, dcId)

	resp, err := c.CloudApiV6Services.Servers().Delete(dcId, serverId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Server successfully deleted")
	return nil
}

func RunServerStart(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "start server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Server is starting... ")

	resp, err := c.CloudApiV6Services.Servers().Start(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Server successfully started")
	return nil
}

func RunServerStop(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "stop server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Server is stopping... ")

	resp, err := c.CloudApiV6Services.Servers().Stop(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Server successfully stopped")
	return nil
}

func RunServerSuspend(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "suspend cube server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Server is Suspending... ")

	resp, err := c.CloudApiV6Services.Servers().Suspend(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Server successfully suspended")
	return nil
}

func RunServerReboot(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "reboot server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Server is rebooting... ")

	resp, err := c.CloudApiV6Services.Servers().Reboot(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Server successfully rebooted")
	return nil
}

func RunServerResume(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "resume cube server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Server is resuming... ")

	resp, err := c.CloudApiV6Services.Servers().Resume(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Server successfully resumed")
	return nil
}

func getUpdateServerInfo(c *core.CommandConfig) (*resources.ServerProperties, error) {
	input := ionoscloud.ServerProperties{}

	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		input.SetName(name)

		c.Verbose("Property name set: %v ", name)
	}

	if c.Flags().Changed(constants.FlagCpuFamily) {
		cpuFamily := c.Flags().String(constants.FlagCpuFamily)
		input.SetCpuFamily(cpuFamily)

		c.Verbose("Property CpuFamily set: %v ", cpuFamily)
	}

	if c.Flags().Changed(constants.FlagNICMultiQueue) {
		input.SetNicMultiQueue(c.Flags().Bool(constants.FlagNICMultiQueue))

		c.Verbose("Property NicMultiQueue set: %v ", c.Flags().Bool(constants.FlagNICMultiQueue))
	}

	if c.Flags().Changed(constants.FlagAvailabilityZone) {
		availabilityZone := c.Flags().String(constants.FlagAvailabilityZone)
		input.SetAvailabilityZone(availabilityZone)

		c.Verbose("Property AvailabilityZone set: %v ", availabilityZone)
	}

	if c.Flags().Changed(constants.FlagCores) {
		cores := c.Flags().Int32(constants.FlagCores)
		input.SetCores(cores)

		c.Verbose("Property Cores set: %v ", cores)
	}

	if c.Flags().Changed(cloudapiv6.ArgVolumeId) {
		volumeId := c.Flags().String(cloudapiv6.ArgVolumeId)
		input.SetBootVolume(ionoscloud.ResourceReference{
			Id: &volumeId,
		})

		c.Verbose("Property BootVolume set: %v ", volumeId)
	}

	if c.Flags().Changed(cloudapiv6.ArgCdromId) {
		cdromId := c.Flags().String(cloudapiv6.ArgCdromId)
		input.SetBootCdrom(ionoscloud.ResourceReference{
			Id: &cdromId,
		})

		c.Verbose("Property BootCdrom set: %v ", cdromId)
	}

	if c.Flags().Changed(constants.FlagRam) {
		size, err := utils2.ConvertSize(
			c.Flags().String(constants.FlagRam),
			utils2.MegaBytes,
		)
		if err != nil {
			return nil, err
		}
		if size < 0 || size > math.MaxInt32 {
			return nil, fmt.Errorf("RAM size %d is out of allowed int32 range [0-%d]", size, math.MaxInt32)
		}
		input.SetRam(int32(size))

		c.Verbose("Property Ram set: %vMB ", int32(size))
	}

	return &resources.ServerProperties{
		ServerProperties: input,
	}, nil
}

func getNewServer(c *core.CommandConfig) (*resources.Server, error) {
	input := resources.ServerProperties{}

	serverType := c.Flags().String(constants.FlagType)
	availabilityZone := c.Flags().String(constants.FlagAvailabilityZone)
	name := c.Flags().String(cloudapiv6.ArgName)

	input.SetType(serverType)
	input.SetAvailabilityZone(availabilityZone)
	input.SetName(name)

	// Confidential VMs derive cores + CPU family from the boot image; leave both unset.
	confidential := c.Flags().Bool(constants.FlagConfidential)

	if c.Flags().Changed(constants.FlagNICMultiQueue) {
		input.SetNicMultiQueue(c.Flags().Bool(constants.FlagNICMultiQueue))

		c.Verbose("Property NicMultiQueue set: %v ", c.Flags().Bool(constants.FlagNICMultiQueue))
	}

	c.Verbose("Property Type set: %v", serverType)
	c.Verbose("Property AvailabilityZone set: %v", availabilityZone)
	c.Verbose("Property Name set: %v", name)

	// GPU Server Properties
	if c.Flags().String(constants.FlagType) == serverGPUType {
		input.ServerProperties.CpuFamily = nil // it automatically selects the correct CPU Family

		if !input.HasName() {
			input.SetName("Unnamed GPU Server")
		}

		if c.Flags().Changed(cloudapiv6.ArgTemplateId) {
			templateUuid := c.Flags().String(cloudapiv6.ArgTemplateId)
			input.SetTemplateUuid(templateUuid)

			c.Verbose("Property TemplateUuid set: %v", templateUuid)
		}
	}

	// CUBE Server Properties
	if c.Flags().String(constants.FlagType) == serverCubeType {
		input.ServerProperties.CpuFamily = nil
		if c.Flags().Changed(constants.FlagCpuFamily) {
			// NOTE 19.07.2023:
			// In the past, all CUBE servers had to have "INTEL_SKYLAKE" as a CPU Family.
			// As such, INTEL_SKYLAKE was hardcoded as the CpuFamily field.
			//
			// However, something changed on the API side, and started throwing errs if Cpu Family was set:
			// `[VDC-5-1921] The attribute 'cpuFamily' must not be provided for Cube servers.`
			//
			// I will allow the user to modify this field, but only if the flag is explicitly set,
			// in case the API changes back to its old state in the future

			input.SetCpuFamily(c.Flags().String(constants.FlagCpuFamily))
		}
		if !input.HasName() {
			input.SetName("Unnamed Cube")
		}
		if c.Flags().Changed(cloudapiv6.ArgTemplateId) {
			templateUuid := c.Flags().String(cloudapiv6.ArgTemplateId)
			input.SetTemplateUuid(templateUuid)

			c.Verbose("Property TemplateUuid set: %v", templateUuid)
		}
	}

	// ENTERPRISE Server Properties
	if c.Flags().String(constants.FlagType) == serverEnterpriseType {
		// For Confidential VMs the CPU family is derived from the image (launch-config vcpu-model);
		// leave it unset so the API resolves it. Otherwise use the flag value or the location default.
		if !confidential {
			if c.Flags().Changed(constants.FlagCpuFamily) &&
				c.Flags().String(constants.FlagCpuFamily) != cloudapiv6.DefaultServerCPUFamily {
				input.SetCpuFamily(c.Flags().String(constants.FlagCpuFamily))
			} else {
				cpuFamily, err := DefaultCpuFamily(c)
				if err != nil {
					return nil, err
				}

				input.SetCpuFamily(cpuFamily)
			}
		}

		if !input.HasName() {
			input.SetName("Unnamed Server")
		}

		// Cores are derived from the image (launch-config vcpu-count) for Confidential VMs.
		if !confidential && c.Flags().Changed(constants.FlagCores) {
			cores := c.Flags().Int32(constants.FlagCores)
			input.SetCores(cores)

			c.Verbose("Property Cores set: %v", cores)
		}

		if c.Flags().Changed(constants.FlagRam) {
			size, err := utils2.ConvertSize(
				c.Flags().String(constants.FlagRam),
				utils2.MegaBytes,
			)
			if err != nil {
				return nil, err
			}

			if size < 0 || size > math.MaxInt32 {
				return nil, fmt.Errorf("RAM size %d is out of allowed int32 range [0-%d]", size, math.MaxInt32)
			}
			input.SetRam(int32(size))

			c.Verbose("Property Ram set: %vMB", int32(size))
		}
	}

	if c.Flags().String(constants.FlagType) == serverVCPUType {
		input.ServerProperties.CpuFamily = nil
		if c.Flags().Changed(constants.FlagCpuFamily) {
			input.SetCpuFamily(c.Flags().String(constants.FlagCpuFamily))
		}

		if !input.HasName() {
			input.SetName("Unnamed VCPU")
		}

		if c.Flags().Changed(constants.FlagCores) {
			cores := c.Flags().Int32(constants.FlagCores)
			input.SetCores(cores)

			c.Verbose("Property Cores set: %v", cores)
		}
		if c.Flags().Changed(constants.FlagRam) {
			size, err := utils2.ConvertSize(
				c.Flags().String(constants.FlagRam),
				utils2.MegaBytes,
			)
			if err != nil {
				return nil, err
			}

			if size < 0 || size > math.MaxInt32 {
				return nil, fmt.Errorf("RAM size %d is out of allowed int32 range [0-%d]", size, math.MaxInt32)
			}
			input.SetRam(int32(size))

			c.Verbose("Property Ram set: %vMB", int32(size))
		}
	}

	return &resources.Server{
		Server: ionoscloud.Server{
			Properties: &input.ServerProperties,
		},
	}, nil
}

func getNewDAS(c *core.CommandConfig) (*resources.Volume, error) {
	volumeProper := resources.VolumeProperties{}

	serverType := c.Flags().String(constants.FlagType)
	if serverType == serverCubeType {
		volumeProper.SetType("DAS")
	}

	// Confidential boot volume: a normal sized volume (not template-based DAS) built from the
	// confidential image. Set its storage type and size from the dedicated flags.
	if c.Flags().Bool(constants.FlagConfidential) {
		volumeProper.SetType(c.Flags().String(constants.FlagStorageType))
		size, err := utils2.ConvertSize(
			c.Flags().String(cloudapiv6.ArgSize),
			utils2.GigaBytes,
		)
		if err != nil {
			return nil, err
		}
		volumeProper.SetSize(float32(size))
	}

	volumeProper.SetName(c.Flags().String(cloudapiv6.ArgVolumeName))
	volumeProper.SetBus(c.Flags().String(cloudapiv6.ArgBus))

	if (!c.Flags().Changed(cloudapiv6.ArgImageId) &&
		!c.Flags().Changed(cloudapiv6.ArgImageAlias)) ||
		c.Flags().Changed(cloudapiv6.ArgLicenceType) {
		volumeProper.SetLicenceType(c.Flags().String(cloudapiv6.ArgLicenceType))
	}

	if c.Flags().Changed(cloudapiv6.ArgImageId) {
		volumeProper.SetImage(c.Flags().String(cloudapiv6.ArgImageId))
	}

	if c.Flags().Changed(cloudapiv6.ArgImageAlias) {
		volumeProper.SetImageAlias(c.Flags().String(cloudapiv6.ArgImageAlias))
	}

	if c.Flags().Changed(constants.ArgPassword) {
		volumeProper.SetImagePassword(c.Flags().String(constants.ArgPassword))
	}

	if c.Flags().Changed(cloudapiv6.ArgSshKeyPaths) {
		sshKeysPaths := c.Flags().StringSlice(cloudapiv6.ArgSshKeyPaths)

		c.Verbose("SSH Key Paths: %v", sshKeysPaths)

		sshKeys, err := helpers.GetSshKeysFromPaths(sshKeysPaths)
		if err != nil {
			return nil, err
		}

		volumeProper.SetSshKeys(sshKeys)

		c.Verbose("Property SshKeys set")
	}

	return &resources.Volume{
		Volume: ionoscloud.Volume{
			Properties: &volumeProper.VolumeProperties,
		},
	}, nil
}

func DeleteAllServers(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Getting Servers...")

	servers, resp, err := c.CloudApiV6Services.Servers().List(dcId)
	if err != nil {
		return err
	}

	serversItems, ok := servers.GetItemsOk()
	if !ok || serversItems == nil {
		return fmt.Errorf("could not get items of Servers")
	}

	if len(*serversItems) <= 0 {
		return fmt.Errorf("no Servers found")
	}

	var multiErr error
	for _, server := range *serversItems {
		id := server.GetId()
		name := server.Properties.Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Server with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Servers().Delete(dcId, *id)
		if resp != nil && request.GetId(resp) != "" {
			c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func DefaultCpuFamily(c *core.CommandConfig) (string, error) {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)

	dc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(), dcId).Execute()
	if err != nil {
		return "", err
	}

	if dc.Properties == nil {
		return "", fmt.Errorf("could not retrieve Datacenter Properties")
	}

	if dc.Properties.CpuArchitecture == nil {
		return "", errors.New("could not retrieve CpuArchitecture")
	}

	cpuArch := (*dc.Properties.CpuArchitecture)[0]

	if cpuArch.CpuFamily == nil {
		return "", errors.New("could not retrieve CpuFamily")
	}

	return *cpuArch.CpuFamily, nil
}

func getRequiredFlagsByServerType(serverType string) ([]string, error) {
	baseRequired := []string{cloudapiv6.ArgDataCenterId}

	switch serverType {
	case serverEnterpriseType:
		return append(baseRequired, constants.FlagCores, constants.FlagRam), nil
	case serverVCPUType:
		return append(baseRequired, constants.FlagType, constants.FlagCores, constants.FlagRam), nil
	case serverGPUType:
		return append(baseRequired, constants.FlagType, cloudapiv6.ArgTemplateId), nil
	case serverCubeType:
		return append(baseRequired, constants.FlagType, cloudapiv6.ArgTemplateId), nil
	default:
		return nil, fmt.Errorf("unknown server type %s (valid: ENTERPRISE | VCPU | CUBE | GPU)", serverType)
	}
}
