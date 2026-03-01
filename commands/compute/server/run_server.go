package server

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
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
	serverType := viper.GetString(core.GetFlagName(c.NS, constants.FlagType))
	requiredFlags, err := getRequiredFlagsByServerType(serverType)
	if err != nil {
		return err
	}

	if err = core.CheckRequiredFlags(c.Command, c.NS, requiredFlags...); err != nil {
		return fmt.Errorf("missing %s flags: %w", serverType, err)
	}

	// CUBE Attached Volume promotion logic (--promote-volume)
	// --promote-volume requires --wait-for-state and --type CUBE/GPU to be set
	c.Command.Command.MarkFlagsRequiredTogether(constants.FlagPromoteVolume, constants.ArgWaitForState)
	if viper.GetBool(core.GetFlagName(c.NS, constants.FlagPromoteVolume)) {
		serverType := viper.GetString(core.GetFlagName(c.NS, constants.FlagType))
		if serverType != serverCubeType && serverType != serverGPUType {
			return fmt.Errorf("--%s can only be used with --%s %s or %s",
				constants.FlagPromoteVolume, constants.FlagType, serverCubeType, serverGPUType)
		}
	}

	imageIdFlag := core.GetFlagName(c.NS, cloudapiv6.ArgImageId)
	imageAliasFlag := core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)

	// Check if image ID or alias is set
	if viper.IsSet(imageIdFlag) || viper.IsSet(imageAliasFlag) {
		imageRequiredFlags := make([][]string, 0)

		if viper.IsSet(imageAliasFlag) {
			// Handle public image alias
			imageRequiredFlags = append(imageRequiredFlags,
				append(requiredFlags, cloudapiv6.ArgImageAlias, cloudapiv6.ArgPassword),
				append(requiredFlags, cloudapiv6.ArgImageAlias, cloudapiv6.ArgSshKeyPaths),
			)
		} else if viper.IsSet(imageIdFlag) {
			// Check if the image ID corresponds to an image or snapshot
			img, _, imgErr := client.Must().CloudClient.ImagesApi.ImagesFindById(context.Background(),
				viper.GetString(imageIdFlag)).Execute()
			if imgErr != nil {
				// Try to fetch as a snapshot if image fetch fails
				_, _, snapshotErr := client.Must().CloudClient.SnapshotsApi.SnapshotsFindById(context.Background(),
					viper.GetString(imageIdFlag)).Execute()
				if snapshotErr != nil {
					return fmt.Errorf("failed getting image or snapshot %s: %w", viper.GetString(imageIdFlag), imgErr)
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
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput(
		"*.items", jsonpaths.Server, allServers, tabheaders.GetHeaders(allServerCols, defaultServerCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunServerList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunServerListAll(c)
	}

	servers, resp, err := c.CloudApiV6Services.Servers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Server, servers.Servers,
		tabheaders.GetHeaders(allServerCols, defaultServerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunServerGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Server with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))))

	if err := waitfor.WaitForState(c, waiter.ServerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))); err != nil {
		return err
	}
	svr, resp, err := c.CloudApiV6Services.Servers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Server, svr.Server,
		tabheaders.GetHeaders(allServerCols, defaultServerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunServerCreate(c *core.CommandConfig) error {
	input, err := getNewServer(c)
	if err != nil {
		return err
	}

	// If Server is of type CUBE, it will create an attached Volume
	if viper.GetString(core.GetFlagName(c.NS, constants.FlagType)) == serverCubeType {
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
	if viper.GetString(core.GetFlagName(c.NS, constants.FlagType)) == serverGPUType {
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

	svr, resp, err := c.CloudApiV6Services.Servers().Create(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), *input)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if id, ok := svr.GetIdOk(); ok && id != nil {
			if err = waitfor.WaitForState(c, waiter.ServerStateInterrogator, *id); err != nil {
				return err
			}

			if svr, _, err = c.CloudApiV6Services.Servers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
				*id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new server id")
		}
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.FlagPromoteVolume)) {
		// Promote the attached Volume to Boot Volume

		if svr.Server.Entities == nil || svr.Server.Entities.Volumes == nil || svr.Server.Entities.Volumes.Items == nil || len(*svr.Server.Entities.Volumes.Items) == 0 {
			return errors.New("no attached volumes found to promote to boot volume")
		}

		attachedDas := (*svr.Server.Entities.Volumes.Items)[0]
		bootVolume := ionoscloud.ResourceReference{
			Id: attachedDas.Id,
		}
		updatedServer, _, err := client.Must().CloudClient.ServersApi.DatacentersServersPatch(context.Background(),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), *svr.Id).
			Server(ionoscloud.ServerProperties{BootVolume: &bootVolume}).Execute()
		if err != nil {
			return fmt.Errorf("error promoting attached volume to boot volume: %w", err)
		}
		svr.Server = updatedServer
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Server, svr.Server,
		tabheaders.GetHeaders(allServerCols, defaultServerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunServerUpdate(c *core.CommandConfig) error {
	input, err := getUpdateServerInfo(c)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Updating Server with ID: %v in Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))))

	svr, resp, err := c.CloudApiV6Services.Servers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		*input,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if err = waitfor.WaitForState(c, waiter.ServerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))); err != nil {
			return err
		}

		if svr, _, err = c.CloudApiV6Services.Servers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))); err != nil {
			return err
		}
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Server, svr.Server,
		tabheaders.GetHeaders(allServerCols, defaultServerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunServerDelete(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllServers(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Starting deleting Server with id: %v from datacenter with id: %v... ", serverId, dcId))

	resp, err := c.CloudApiV6Services.Servers().Delete(dcId, serverId)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Server successfully deleted"))
	return nil
}

func RunServerStart(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "start server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Server is starting... "))

	resp, err := c.CloudApiV6Services.Servers().Start(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Server successfully started"))
	return nil
}

func RunServerStop(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "stop server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Server is stopping... "))

	resp, err := c.CloudApiV6Services.Servers().Stop(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Server successfully stopped"))
	return nil
}

func RunServerSuspend(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "suspend cube server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Server is Suspending... "))

	resp, err := c.CloudApiV6Services.Servers().Suspend(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Server successfully suspended"))
	return nil
}

func RunServerReboot(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "reboot server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Server is rebooting... "))

	resp, err := c.CloudApiV6Services.Servers().Reboot(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Server successfully rebooted"))
	return nil
}

func RunServerResume(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "resume cube server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Server is resuming... "))

	resp, err := c.CloudApiV6Services.Servers().Resume(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Server successfully resumed"))
	return nil
}

func getUpdateServerInfo(c *core.CommandConfig) (*resources.ServerProperties, error) {
	input := ionoscloud.ServerProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property name set: %v ", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCpuFamily)) {
		cpuFamily := viper.GetString(core.GetFlagName(c.NS, constants.FlagCpuFamily))
		input.SetCpuFamily(cpuFamily)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property CpuFamily set: %v ", cpuFamily))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagNICMultiQueue); viper.IsSet(fn) {
		input.SetNicMultiQueue(viper.GetBool(fn))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property NicMultiQueue set: %v ", viper.GetBool(fn)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagAvailabilityZone)) {
		availabilityZone := viper.GetString(core.GetFlagName(c.NS, constants.FlagAvailabilityZone))
		input.SetAvailabilityZone(availabilityZone)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property AvailabilityZone set: %v ", availabilityZone))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
		cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
		input.SetCores(cores)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Cores set: %v ", cores))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)) {
		volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
		input.SetBootVolume(ionoscloud.ResourceReference{
			Id: &volumeId,
		})

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property BootVolume set: %v ", volumeId))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId)) {
		cdromId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId))
		input.SetBootCdrom(ionoscloud.ResourceReference{
			Id: &cdromId,
		})

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property BootCdrom set: %v ", cdromId))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
		size, err := utils2.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)),
			utils2.MegaBytes,
		)
		if err != nil {
			return nil, err
		}
		if size < 0 || size > math.MaxInt32 {
			return nil, fmt.Errorf("RAM size %d is out of allowed int32 range [0-%d]", size, math.MaxInt32)
		}
		input.SetRam(int32(size))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Ram set: %vMB ", int32(size)))
	}

	return &resources.ServerProperties{
		ServerProperties: input,
	}, nil
}

func getNewServer(c *core.CommandConfig) (*resources.Server, error) {
	input := resources.ServerProperties{}

	serverType := viper.GetString(core.GetFlagName(c.NS, constants.FlagType))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, constants.FlagAvailabilityZone))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))

	input.SetType(serverType)
	input.SetAvailabilityZone(availabilityZone)
	input.SetName(name)

	if fn := core.GetFlagName(c.NS, constants.FlagNICMultiQueue); viper.IsSet(fn) {
		input.SetNicMultiQueue(viper.GetBool(fn))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property NicMultiQueue set: %v ", viper.GetBool(fn)))
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Type set: %v", serverType))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property AvailabilityZone set: %v", availabilityZone))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))

	// GPU Server Properties
	if viper.GetString(core.GetFlagName(c.NS, constants.FlagType)) == serverGPUType {
		input.ServerProperties.CpuFamily = nil // it automatically selects the correct CPU Family

		if !input.HasName() {
			input.SetName("Unnamed GPU Server")
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId)) {
			templateUuid := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId))
			input.SetTemplateUuid(templateUuid)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property TemplateUuid set: %v", templateUuid))
		}
	}

	// CUBE Server Properties
	if viper.GetString(core.GetFlagName(c.NS, constants.FlagType)) == serverCubeType {
		input.ServerProperties.CpuFamily = nil
		if fn := core.GetFlagName(c.NS, constants.FlagCpuFamily); viper.IsSet(fn) {
			// NOTE 19.07.2023:
			// In the past, all CUBE servers had to have "INTEL_SKYLAKE" as a CPU Family.
			// As such, INTEL_SKYLAKE was hardcoded as the CpuFamily field.
			//
			// However, something changed on the API side, and started throwing errs if Cpu Family was set:
			// `[VDC-5-1921] The attribute 'cpuFamily' must not be provided for Cube servers.`
			//
			// I will allow the user to modify this field, but only if the flag is explicitly set,
			// in case the API changes back to its old state in the future

			input.SetCpuFamily(viper.GetString(fn))
		}
		if !input.HasName() {
			input.SetName("Unnamed Cube")
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId)) {
			templateUuid := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId))
			input.SetTemplateUuid(templateUuid)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property TemplateUuid set: %v", templateUuid))
		}
	}

	// ENTERPRISE Server Properties
	if viper.GetString(core.GetFlagName(c.NS, constants.FlagType)) == serverEnterpriseType {
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCpuFamily)) &&
			viper.GetString(core.GetFlagName(c.NS, constants.FlagCpuFamily)) != cloudapiv6.DefaultServerCPUFamily {
			input.SetCpuFamily(viper.GetString(core.GetFlagName(c.NS, constants.FlagCpuFamily)))
		} else {
			cpuFamily, err := DefaultCpuFamily(c)
			if err != nil {
				return nil, err
			}

			input.SetCpuFamily(cpuFamily)
		}

		if !input.HasName() {
			input.SetName("Unnamed Server")
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
			cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
			input.SetCores(cores)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Cores set: %v", cores))
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
			size, err := utils2.ConvertSize(
				viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)),
				utils2.MegaBytes,
			)
			if err != nil {
				return nil, err
			}

			if size < 0 || size > math.MaxInt32 {
				return nil, fmt.Errorf("RAM size %d is out of allowed int32 range [0-%d]", size, math.MaxInt32)
			}
			input.SetRam(int32(size))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Ram set: %vMB", int32(size)))
		}
	}

	if viper.GetString(core.GetFlagName(c.NS, constants.FlagType)) == serverVCPUType {
		input.ServerProperties.CpuFamily = nil
		if fn := core.GetFlagName(c.NS, constants.FlagCpuFamily); viper.IsSet(fn) {
			input.SetCpuFamily(viper.GetString(fn))
		}

		if !input.HasName() {
			input.SetName("Unnamed VCPU")
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
			cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
			input.SetCores(cores)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Cores set: %v", cores))
		}
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
			size, err := utils2.ConvertSize(
				viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)),
				utils2.MegaBytes,
			)
			if err != nil {
				return nil, err
			}

			if size < 0 || size > math.MaxInt32 {
				return nil, fmt.Errorf("RAM size %d is out of allowed int32 range [0-%d]", size, math.MaxInt32)
			}
			input.SetRam(int32(size))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Ram set: %vMB", int32(size)))
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

	serverType := viper.GetString(core.GetFlagName(c.NS, constants.FlagType))
	if serverType == serverCubeType {
		volumeProper.SetType("DAS")
	}
	volumeProper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeName)))
	volumeProper.SetBus(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBus)))

	if (!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) &&
		!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType)) {
		volumeProper.SetLicenceType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) {
		volumeProper.SetImage(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)) {
		volumeProper.SetImageAlias(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgPassword)) {
		volumeProper.SetImagePassword(viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths)) {
		sshKeysPaths := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("SSH Key Paths: %v", sshKeysPaths))

		sshKeys, err := helpers.GetSshKeysFromPaths(sshKeysPaths)
		if err != nil {
			return nil, err
		}

		volumeProper.SetSshKeys(sshKeys)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property SshKeys set"))
	}

	return &resources.Volume{
		Volume: ionoscloud.Volume{
			Properties: &volumeProper.VolumeProperties,
		},
	}, nil
}

func DeleteAllServers(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Servers..."))

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
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func DefaultCpuFamily(c *core.CommandConfig) (string, error) {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

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
