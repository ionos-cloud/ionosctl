package snapshot

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

func PreRunSnapshotId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgSnapshotId)
}

func PreRunSnapshotDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgSnapshotId},
		[]string{cloudapiv6.ArgAll},
	)
}

func PreRunSnapshotIdDcIdVolumeId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgVolumeId, cloudapiv6.ArgSnapshotId)
}

func RunSnapshotList(c *core.CommandConfig) error {
	ss, resp, err := c.CloudApiV6Services.Snapshots().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allSnapshotCols).Prefix("items").Print(ss.Snapshots)
}

func RunSnapshotGet(c *core.CommandConfig) error {
	c.Verbose("Snapshot with id: %v is getting...", c.Flags().String(cloudapiv6.ArgSnapshotId))

	s, resp, err := c.CloudApiV6Services.Snapshots().Get(c.Flags().String(cloudapiv6.ArgSnapshotId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allSnapshotCols).Print(s.Snapshot)
}

func RunSnapshotCreate(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	volumeId := c.Flags().String(cloudapiv6.ArgVolumeId)
	name := c.Flags().String(cloudapiv6.ArgName)
	description := c.Flags().String(cloudapiv6.ArgDescription)
	licenseType := c.Flags().String(cloudapiv6.ArgLicenceType)
	secAuthProtection := c.Flags().Bool(cloudapiv6.ArgSecAuthProtection)

	c.Verbose("Properties set for creating the Snapshot: DatacenterId: %v, VolumeId: %v, Name: %v, Description: %v, "+
		"LicenseType: %v, SecAuthProtection: %v", dcId, volumeId, name, description, licenseType, secAuthProtection)

	s, resp, err := c.CloudApiV6Services.Snapshots().Create(dcId, volumeId, name, description, licenseType, secAuthProtection)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allSnapshotCols).Print(s.Snapshot)
}

func RunSnapshotUpdate(c *core.CommandConfig) error {
	c.Verbose("Updating Snapshot with id: %v...", c.Flags().String(cloudapiv6.ArgSnapshotId))

	s, resp, err := c.CloudApiV6Services.Snapshots().Update(c.Flags().String(cloudapiv6.ArgSnapshotId),
		getSnapshotPropertiesSet(c))
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allSnapshotCols).Print(s.Snapshot)
}

func RunSnapshotRestore(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "restore snapshot", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Snapshot with id: %v is restoring...", c.Flags().String(cloudapiv6.ArgSnapshotId))

	resp, err := c.CloudApiV6Services.Snapshots().Restore(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgVolumeId),
		c.Flags().String(cloudapiv6.ArgSnapshotId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Snapshot successfully restored")
	return nil
}

func RunSnapshotDelete(c *core.CommandConfig) error {
	snapshotId := c.Flags().String(cloudapiv6.ArgSnapshotId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DeleteAllSnapshots(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete snapshot", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Snapshot with id: %v...", snapshotId)

	resp, err := c.CloudApiV6Services.Snapshots().Delete(snapshotId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Snapshot successfully deleted")
	return nil
}

func getSnapshotPropertiesSet(c *core.CommandConfig) resources.SnapshotProperties {
	input := resources.SnapshotProperties{}

	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if c.Flags().Changed(cloudapiv6.ArgDescription) {
		description := c.Flags().String(cloudapiv6.ArgDescription)
		input.SetDescription(description)

		c.Verbose("Property Description set: %v", description)
	}

	if c.Flags().Changed(cloudapiv6.ArgLicenceType) {
		licenceType := c.Flags().String(cloudapiv6.ArgLicenceType)
		input.SetLicenceType(licenceType)

		c.Verbose("Property LicenceType set: %v", licenceType)
	}

	if c.Flags().Changed(cloudapiv6.ArgCpuHotPlug) {
		cpuHotPlug := c.Flags().Bool(cloudapiv6.ArgCpuHotPlug)
		input.SetCpuHotPlug(cpuHotPlug)

		c.Verbose("Property CpuHotPlug set: %v", cpuHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgCpuHotUnplug) {
		cpuHotUnplug := c.Flags().Bool(cloudapiv6.ArgCpuHotUnplug)
		input.SetCpuHotUnplug(cpuHotUnplug)

		c.Verbose("Property CpuHotUnplug set: %v", cpuHotUnplug)
	}

	if c.Flags().Changed(cloudapiv6.ArgRamHotPlug) {
		ramHotPlug := c.Flags().Bool(cloudapiv6.ArgRamHotPlug)
		input.SetRamHotPlug(ramHotPlug)

		c.Verbose("Property RamHotPlug set: %v", ramHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgRamHotUnplug) {
		ramHotUnplug := c.Flags().Bool(cloudapiv6.ArgRamHotUnplug)
		input.SetRamHotUnplug(ramHotUnplug)

		c.Verbose("Property RamHotUnplug set: %v", ramHotUnplug)
	}

	if c.Flags().Changed(cloudapiv6.ArgNicHotPlug) {
		nicHotPlug := c.Flags().Bool(cloudapiv6.ArgNicHotPlug)
		input.SetNicHotPlug(nicHotPlug)

		c.Verbose("Property NicHotPlug set: %v", nicHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgNicHotUnplug) {
		nicHotUnplug := c.Flags().Bool(cloudapiv6.ArgNicHotUnplug)
		input.SetNicHotUnplug(nicHotUnplug)

		c.Verbose("Property nicHotUnplug set: %v", nicHotUnplug)
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

	if c.Flags().Changed(cloudapiv6.ArgDiscScsiHotPlug) {
		discScsiHotPlug := c.Flags().Bool(cloudapiv6.ArgDiscScsiHotPlug)
		input.SetDiscScsiHotPlug(discScsiHotPlug)

		c.Verbose("Property DiscScsiHotPlug set: %v", discScsiHotPlug)
	}

	if c.Flags().Changed(cloudapiv6.ArgDiscScsiHotUnplug) {
		discScsiHotUnplug := c.Flags().Bool(cloudapiv6.ArgDiscScsiHotUnplug)
		input.SetDiscScsiHotUnplug(discScsiHotUnplug)

		c.Verbose("Property DiscScsiHotUnplug set: %v", discScsiHotUnplug)
	}

	if c.Flags().Changed(cloudapiv6.ArgSecAuthProtection) {
		secAuthProtection := c.Flags().Bool(cloudapiv6.ArgSecAuthProtection)
		input.SetSecAuthProtection(secAuthProtection)

		c.Verbose("Property SecAuthProtection set: %v", secAuthProtection)
	}

	return input
}

func DeleteAllSnapshots(c *core.CommandConfig) error {
	c.Verbose("Getting Snapshots...")

	snapshots, resp, err := c.CloudApiV6Services.Snapshots().List()
	if err != nil {
		return err
	}

	snapshotsItems, ok := snapshots.GetItemsOk()
	if !ok || snapshotsItems == nil {
		return fmt.Errorf("could not get items of Snapshots")
	}

	if len(*snapshotsItems) <= 0 {
		return fmt.Errorf("no Snapshots found")
	}

	var multiErr error
	for _, snapshot := range *snapshotsItems {
		id := snapshot.GetId()
		name := snapshot.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Snapshot with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Snapshots().Delete(*id)
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

func PreRunDcVolumeIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgVolumeId)
}
