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
	c.Verbose("Snapshot with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))

	s, resp, err := c.CloudApiV6Services.Snapshots().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allSnapshotCols).Print(s.Snapshot)
}

func RunSnapshotCreate(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))
	licenseType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType))
	secAuthProtection := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgSecAuthProtection))

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
	c.Verbose("Updating Snapshot with id: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))

	s, resp, err := c.CloudApiV6Services.Snapshots().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)),
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

	c.Verbose("Snapshot with id: %v is restoring...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))

	resp, err := c.CloudApiV6Services.Snapshots().Restore(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)),
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
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
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

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDescription)) {
		description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))
		input.SetDescription(description)

		c.Verbose("Property Description set: %v", description)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType)) {
		licenceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType))
		input.SetLicenceType(licenceType)

		c.Verbose("Property LicenceType set: %v", licenceType)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug))
		input.SetCpuHotPlug(cpuHotPlug)

		c.Verbose("Property CpuHotPlug set: %v", cpuHotPlug)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotUnplug)) {
		cpuHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotUnplug))
		input.SetCpuHotUnplug(cpuHotUnplug)

		c.Verbose("Property CpuHotUnplug set: %v", cpuHotUnplug)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug))
		input.SetRamHotPlug(ramHotPlug)

		c.Verbose("Property RamHotPlug set: %v", ramHotPlug)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotUnplug)) {
		ramHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotUnplug))
		input.SetRamHotUnplug(ramHotUnplug)

		c.Verbose("Property RamHotUnplug set: %v", ramHotUnplug)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug))
		input.SetNicHotPlug(nicHotPlug)

		c.Verbose("Property NicHotPlug set: %v", nicHotPlug)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug))
		input.SetNicHotUnplug(nicHotUnplug)

		c.Verbose("Property nicHotUnplug set: %v", nicHotUnplug)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug))
		input.SetDiscVirtioHotPlug(discVirtioHotPlug)

		c.Verbose("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug))
		input.SetDiscVirtioHotUnplug(discVirtioHotUnplug)

		c.Verbose("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscScsiHotPlug)) {
		discScsiHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscScsiHotPlug))
		input.SetDiscScsiHotPlug(discScsiHotPlug)

		c.Verbose("Property DiscScsiHotPlug set: %v", discScsiHotPlug)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscScsiHotUnplug)) {
		discScsiHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscScsiHotUnplug))
		input.SetDiscScsiHotUnplug(discScsiHotUnplug)

		c.Verbose("Property DiscScsiHotUnplug set: %v", discScsiHotUnplug)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSecAuthProtection)) {
		secAuthProtection := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgSecAuthProtection))
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
