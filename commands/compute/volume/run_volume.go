package volume

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
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

	if c.Flags().Changed(cloudapiv6.ArgImageId) {
		// Define required flags for private images
		setRequiredFlagsPrivateImage := [][]string{
			{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageId},
			{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageAlias},
		}

		img, _, err := client.Must().CloudClient.ImagesApi.ImagesFindById(context.Background(),
			c.Flags().String(cloudapiv6.ArgImageId)).Execute()
		if err != nil {
			// try to fetch it as a snapshot if fails
			_, _, snapshotErr := client.Must().CloudClient.SnapshotsApi.SnapshotsFindById(context.Background(),
				c.Flags().String(cloudapiv6.ArgImageId)).Execute()
			if snapshotErr != nil {
				return fmt.Errorf("failed getting image or snapshot %s: %w", c.Flags().String(cloudapiv6.ArgImageId), err)
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
	if publicImageAsImageId || c.Flags().Changed(cloudapiv6.ArgImageAlias) {
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
		c.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer(allVolumeCols).Prefix("*.items").Print(allVolumes)
}

func RunVolumeList(c *core.CommandConfig) error {
	if c.Flags().Bool(cloudapiv6.ArgAll) {
		return RunVolumeListAll(c)
	}

	c.Verbose("Listing Volumes from Datacenter with ID: %v", c.Flags().String(cloudapiv6.ArgDataCenterId))

	volumes, resp, err := c.CloudApiV6Services.Volumes().List(c.Flags().String(cloudapiv6.ArgDataCenterId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allVolumeCols).Prefix("items").Print(volumes.Volumes)
}

func RunVolumeGet(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, c.Flags().String(cloudapiv6.ArgDataCenterId))
	c.Verbose("Volume with id: %v is getting...", c.Flags().String(cloudapiv6.ArgVolumeId))

	vol, resp, err := c.CloudApiV6Services.Volumes().Get(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgVolumeId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allVolumeCols).Print(vol.Volume)
}

func RunVolumeCreate(c *core.CommandConfig) error {
	input, err := getNewVolume(c)
	if err != nil {
		return err
	}

	vol, resp, err := c.CloudApiV6Services.Volumes().Create(c.Flags().String(cloudapiv6.ArgDataCenterId), *input)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allVolumeCols).Print(vol.Volume)
}

func RunVolumeUpdate(c *core.CommandConfig) error {
	input, err := getVolumeInfo(c)
	if err != nil {
		return err
	}

	vol, resp, err := c.CloudApiV6Services.Volumes().Update(c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgVolumeId), *input)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allVolumeCols).Print(vol.Volume)
}

func RunVolumeDelete(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	volumeId := c.Flags().String(cloudapiv6.ArgVolumeId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DeleteAllVolumes(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete volume", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Volume with id: %v...", volumeId)

	resp, err := c.CloudApiV6Services.Volumes().Delete(dcId, volumeId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Volume successfully deleted")
	return nil
}

func getNewVolume(c *core.CommandConfig) (*resources.Volume, error) {
	proper := resources.VolumeProperties{}

	name := c.Flags().String(cloudapiv6.ArgName)
	bus := c.Flags().String(cloudapiv6.ArgBus)
	volumeType := c.Flags().String(constants.FlagType)
	availabilityZone := c.Flags().String(constants.FlagAvailabilityZone)

	// It will get the default values, if flags not set
	proper.SetName(name)
	proper.SetBus(bus)
	proper.SetType(volumeType)
	proper.SetAvailabilityZone(availabilityZone)

	c.Verbose("Properties set for creating the Volume: Name: %v, Bus: %v, VolumeType: %v, AvailabilityZone: %v",
		name, bus, volumeType, availabilityZone)

	size, err := utils2.ConvertSize(
		c.Flags().String(cloudapiv6.ArgSize),
		utils2.GigaBytes,
	)
	if err != nil {
		return nil, err
	}

	proper.SetSize(float32(size))

	c.Verbose("Property Size set: %vGB", float32(size))

	// Check if flags are set and set options
	if c.Flags().Changed(cloudapiv6.ArgBackupUnitId) {
		backupUnitId := c.Flags().String(cloudapiv6.ArgBackupUnitId)
		proper.SetBackupunitId(backupUnitId)

		c.Verbose("Property BackupUnitId set: %v", backupUnitId)
	}

	if (!c.Flags().Changed(cloudapiv6.ArgImageId) &&
		!c.Flags().Changed(cloudapiv6.ArgImageAlias)) ||
		c.Flags().Changed(cloudapiv6.ArgLicenceType) {
		licenceType := c.Flags().String(cloudapiv6.ArgLicenceType)
		proper.SetLicenceType(licenceType)

		c.Verbose("Property LicenceType set: %v", licenceType)
	}

	if c.Flags().Changed(cloudapiv6.ArgImageId) {
		imageId := c.Flags().String(cloudapiv6.ArgImageId)
		proper.SetImage(imageId)

		c.Verbose("Property Image set: %v", imageId)
	}

	if c.Flags().Changed(cloudapiv6.ArgImageAlias) {
		imageAlias := c.Flags().String(cloudapiv6.ArgImageAlias)
		proper.SetImageAlias(imageAlias)

		c.Verbose("Property ImageAlias set: %v", imageAlias)
	}

	if c.Flags().Changed(cloudapiv6.ArgPassword) {
		imagePassword := c.Flags().String(cloudapiv6.ArgPassword)
		proper.SetImagePassword(imagePassword)

		c.Verbose("Property ImagePassword set")
	}

	if c.Flags().Changed(cloudapiv6.ArgSshKeyPaths) {
		sshKeysPaths := c.Flags().StringSlice(cloudapiv6.ArgSshKeyPaths)

		c.Verbose("SSH Key Paths: %v", sshKeysPaths)

		sshKeys, err := helpers.GetSshKeysFromPaths(sshKeysPaths)
		if err != nil {
			return nil, err
		}

		proper.SetSshKeys(sshKeys)

		c.Verbose("Property SshKeys set")
	}

	if c.Flags().Changed(cloudapiv6.ArgUserData) {
		userData := c.Flags().String(cloudapiv6.ArgUserData)
		proper.SetUserData(userData)

		c.Verbose("Property UserData set: %v", userData)
	}

	if c.Flags().Changed(cloudapiv6.ArgCpuHotPlug) {
		cpuHotPlug := c.Flags().Bool(cloudapiv6.ArgCpuHotPlug)
		proper.SetCpuHotPlug(cpuHotPlug)

		c.Verbose("Property CpuHotPlug set: %v", cpuHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgRamHotPlug) {
		ramHotPlug := c.Flags().Bool(cloudapiv6.ArgRamHotPlug)
		proper.SetRamHotPlug(ramHotPlug)

		c.Verbose("Property RamHotPlug set: %v", ramHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgNicHotPlug) {
		nicHotPlug := c.Flags().Bool(cloudapiv6.ArgNicHotPlug)
		proper.SetNicHotPlug(nicHotPlug)

		c.Verbose("Property NicHotPlug set: %v", nicHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgNicHotUnplug) {
		nicHotUnplug := c.Flags().Bool(cloudapiv6.ArgNicHotUnplug)
		proper.SetNicHotUnplug(nicHotUnplug)

		c.Verbose("Property NicHotUnplug set: %v", nicHotUnplug)
	}

	if c.Flags().Changed(cloudapiv6.ArgDiscVirtioHotPlug) {
		discVirtioHotPlug := c.Flags().Bool(cloudapiv6.ArgDiscVirtioHotPlug)
		proper.SetDiscVirtioHotPlug(discVirtioHotPlug)

		c.Verbose("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgDiscVirtioHotUnplug) {
		discVirtioHotUnplug := c.Flags().Bool(cloudapiv6.ArgDiscVirtioHotUnplug)
		proper.SetDiscVirtioHotUnplug(discVirtioHotUnplug)

		c.Verbose("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug)
	}

	return &resources.Volume{
		Volume: ionoscloud.Volume{
			Properties: &proper.VolumeProperties,
		},
	}, nil
}

func getVolumeInfo(c *core.CommandConfig) (*resources.VolumeProperties, error) {
	input := resources.VolumeProperties{}

	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if c.Flags().Changed(cloudapiv6.ArgBus) {
		bus := c.Flags().String(cloudapiv6.ArgBus)
		input.SetBus(bus)

		c.Verbose("Property Bus set: %v", bus)
	}

	if c.Flags().Changed(cloudapiv6.ArgSize) {
		size, err := utils2.ConvertSize(
			c.Flags().String(cloudapiv6.ArgSize),
			utils2.GigaBytes,
		)
		if err != nil {
			return nil, err
		}

		input.SetSize(float32(size))

		c.Verbose("Property Size set: %vGB", float32(size))
	}

	// Check if flags are set and set options
	if c.Flags().Changed(cloudapiv6.ArgCpuHotPlug) {
		cpuHotPlug := c.Flags().Bool(cloudapiv6.ArgCpuHotPlug)
		input.SetCpuHotPlug(cpuHotPlug)

		c.Verbose("Property CpuHotPlug set: %v", cpuHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgRamHotPlug) {
		ramHotPlug := c.Flags().Bool(cloudapiv6.ArgRamHotPlug)
		input.SetRamHotPlug(ramHotPlug)

		c.Verbose("Property RamHotPlug set: %v", ramHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgNicHotPlug) {
		nicHotPlug := c.Flags().Bool(cloudapiv6.ArgNicHotPlug)
		input.SetNicHotPlug(nicHotPlug)

		c.Verbose("Property NicHotPlug set: %v", nicHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgNicHotUnplug) {
		nicHotUnplug := c.Flags().Bool(cloudapiv6.ArgNicHotUnplug)
		input.SetNicHotUnplug(nicHotUnplug)

		c.Verbose("Property NicHotUnplug set: %v", nicHotUnplug)
	}

	if c.Flags().Changed(cloudapiv6.ArgDiscVirtioHotPlug) {
		discVirtioHotPlug := c.Flags().Bool(cloudapiv6.ArgDiscVirtioHotPlug)
		input.SetDiscVirtioHotPlug(discVirtioHotPlug)

		c.Verbose("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgDiscVirtioHotUnplug) {
		discVirtioHotUnplug := c.Flags().Bool(cloudapiv6.ArgDiscVirtioHotUnplug)
		input.SetDiscVirtioHotUnplug(discVirtioHotUnplug)

		c.Verbose("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug)
	}

	return &input, nil
}

func DeleteAllVolumes(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Getting Volumes...")

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

	c.Msg("Volumes to be deleted:")

	var multiErr error
	for _, volume := range *volumesItems {
		id := volume.GetId()
		name := volume.GetProperties().Name
		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Volume with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Volumes().Delete(dcId, *id)
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
