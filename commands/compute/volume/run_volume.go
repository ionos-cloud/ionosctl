package volume

import (
	"context"
	"errors"
	"fmt"
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

func PreRunVolumeList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	return nil
}

func PreRunDcVolumeIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgVolumeId)
}

func PreRunDcVolumeDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgVolumeId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func PreRunVolumeCreate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId)
	if err != nil {
		return err
	}

	setRequiredFlagsPublicImage := [][]string{
		{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageId, constants.ArgPassword},
		{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageId, cloudapiv6.ArgSshKeyPaths},
		{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageAlias, constants.ArgPassword},
		{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageAlias, cloudapiv6.ArgSshKeyPaths},
	}
	publicImageAsImageId := false

	if fn := core.GetFlagName(c.NS, cloudapiv6.ArgImageId); viper.IsSet(fn) {
		// Define required flags for private images
		setRequiredFlagsPrivateImage := [][]string{
			{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageId},
			{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageAlias},
		}

		img, _, err := client.Must().CloudClient.ImagesApi.ImagesFindById(context.Background(),
			viper.GetString(fn)).Execute()
		if err != nil {
			// try to fetch it as a snapshot if fails
			_, _, snapshotErr := client.Must().CloudClient.SnapshotsApi.SnapshotsFindById(context.Background(),
				viper.GetString(fn)).Execute()
			if snapshotErr != nil {
				return fmt.Errorf("failed getting image or snapshot %s: %w", viper.GetString(fn), err)
			}

			// If a snapshot is found, skip additional checks
			return nil
		}

		if img.Properties == nil || img.Properties.Public == nil || !*img.Properties.Public {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, setRequiredFlagsPrivateImage...)
		}

		publicImageAsImageId = true
	}

	// check public image alias requirements
	if fn := core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias); publicImageAsImageId || viper.IsSet(fn) {
		return core.CheckRequiredFlagsSets(c.Command, c.NS, setRequiredFlagsPublicImage...)
	}

	return nil
}

func RunVolumeListAll(c *core.CommandConfig) error {
	datacenters, _, err := c.CloudApiV6Services.DataCenters().List()
	if err != nil {
		return err
	}

	allDcs := helpers.GetDataCenters(datacenters)
	var allVolumes []ionoscloud.Volumes
	totalTime := time.Duration(0)

	for _, dc := range allDcs {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("could not retrieve Datacenter Id")
		}

		volumes, resp, err := c.CloudApiV6Services.Volumes().List(*dc.GetId())
		if err != nil {
			return err
		}

		allVolumes = append(allVolumes, volumes.Volumes)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput(
		"*.items", jsonpaths.Volume, allVolumes, tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunVolumeList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunVolumeListAll(c)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Listing Volumes from Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))))

	volumes, resp, err := c.CloudApiV6Services.Volumes().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Volume, volumes.Volumes,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunVolumeGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Volume with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))))

	vol, resp, err := c.CloudApiV6Services.Volumes().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Volume, vol.Volume,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunVolumeCreate(c *core.CommandConfig) error {
	input, err := getNewVolume(c)
	if err != nil {
		return err
	}

	vol, resp, err := c.CloudApiV6Services.Volumes().Create(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), *input)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Volume, vol.Volume,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunVolumeUpdate(c *core.CommandConfig) error {
	input, err := getVolumeInfo(c)
	if err != nil {
		return err
	}

	vol, resp, err := c.CloudApiV6Services.Volumes().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)), *input)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Volume, vol.Volume,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunVolumeDelete(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllVolumes(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete volume", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Starting deleting Volume with id: %v...", volumeId))

	resp, err := c.CloudApiV6Services.Volumes().Delete(dcId, volumeId)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Volume successfully deleted"))
	return nil
}

func getNewVolume(c *core.CommandConfig) (*resources.Volume, error) {
	proper := resources.VolumeProperties{}

	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	bus := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBus))
	volumeType := viper.GetString(core.GetFlagName(c.NS, constants.FlagType))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, constants.FlagAvailabilityZone))

	// It will get the default values, if flags not set
	proper.SetName(name)
	proper.SetBus(bus)
	proper.SetType(volumeType)
	proper.SetAvailabilityZone(availabilityZone)

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Properties set for creating the Volume: Name: %v, Bus: %v, VolumeType: %v, AvailabilityZone: %v",
		name, bus, volumeType, availabilityZone))

	size, err := utils2.ConvertSize(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSize)),
		utils2.GigaBytes,
	)
	if err != nil {
		return nil, err
	}

	proper.SetSize(float32(size))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Size set: %vGB", float32(size)))

	// Check if flags are set and set options
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId)) {
		backupUnitId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId))
		proper.SetBackupunitId(backupUnitId)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property BackupUnitId set: %v", backupUnitId))
	}

	if (!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) &&
		!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType)) {
		licenceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType))
		proper.SetLicenceType(licenceType)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property LicenceType set: %v", licenceType))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) {
		imageId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId))
		proper.SetImage(imageId)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Image set: %v", imageId))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)) {
		imageAlias := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))
		proper.SetImageAlias(imageAlias)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property ImageAlias set: %v", imageAlias))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPassword)) {
		imagePassword := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPassword))
		proper.SetImagePassword(imagePassword)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property ImagePassword set"))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths)) {
		sshKeysPaths := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("SSH Key Paths: %v", sshKeysPaths))

		sshKeys, err := helpers.GetSshKeysFromPaths(sshKeysPaths)
		if err != nil {
			return nil, err
		}

		proper.SetSshKeys(sshKeys)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property SshKeys set"))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgUserData)) {
		userData := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserData))
		proper.SetUserData(userData)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property UserData set: %v", userData))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug))
		proper.SetCpuHotPlug(cpuHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property CpuHotPlug set: %v", cpuHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug))
		proper.SetRamHotPlug(ramHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property RamHotPlug set: %v", ramHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug))
		proper.SetNicHotPlug(nicHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property NicHotPlug set: %v", nicHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug))
		proper.SetNicHotUnplug(nicHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property NicHotUnplug set: %v", nicHotUnplug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug))
		proper.SetDiscVirtioHotPlug(discVirtioHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug))
		proper.SetDiscVirtioHotUnplug(discVirtioHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug))
	}

	return &resources.Volume{
		Volume: ionoscloud.Volume{
			Properties: &proper.VolumeProperties,
		},
	}, nil
}

func getVolumeInfo(c *core.CommandConfig) (*resources.VolumeProperties, error) {
	input := resources.VolumeProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgBus)) {
		bus := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBus))
		input.SetBus(bus)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Bus set: %v", bus))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSize)) {
		size, err := utils2.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSize)),
			utils2.GigaBytes,
		)
		if err != nil {
			return nil, err
		}

		input.SetSize(float32(size))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Size set: %vGB", float32(size)))
	}

	// Check if flags are set and set options
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug))
		input.SetCpuHotPlug(cpuHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property CpuHotPlug set: %v", cpuHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug))
		input.SetRamHotPlug(ramHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property RamHotPlug set: %v", ramHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug))
		input.SetNicHotPlug(nicHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property NicHotPlug set: %v", nicHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug))
		input.SetNicHotUnplug(nicHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property NicHotUnplug set: %v", nicHotUnplug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug))
		input.SetDiscVirtioHotPlug(discVirtioHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug))
		input.SetDiscVirtioHotUnplug(discVirtioHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug))
	}

	return &input, nil
}

func DeleteAllVolumes(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Volumes..."))

	volumes, resp, err := c.CloudApiV6Services.Volumes().List(dcId)
	if err != nil {
		return err
	}

	volumesItems, ok := volumes.GetItemsOk()
	if !ok || volumesItems == nil {
		return fmt.Errorf("could not get items of Volumes")
	}

	if len(*volumesItems) <= 0 {
		return fmt.Errorf("no Volumes found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Volumes to be deleted:"))

	var multiErr error
	for _, volume := range *volumesItems {
		id := volume.GetId()
		name := volume.GetProperties().Name
		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Volume with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Volumes().Delete(dcId, *id)
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
