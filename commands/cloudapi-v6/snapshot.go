package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SnapshotCmd() *core.Command {
	ctx := context.TODO()
	snapshotCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "snapshot",
			Aliases:          []string{"ss", "snap"},
			Short:            "Snapshot Operations",
			Long:             "The sub-commands of `ionosctl snapshot` allow you to see information, to create, update, delete Snapshots.",
			TraverseChildren: true,
		},
	}
	globalFlags := snapshotCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultSnapshotCols, printer.ColsMessage(defaultSnapshotCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(snapshotCmd.NS, config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = snapshotCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultSnapshotCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace:  "snapshot",
		Resource:   "snapshot",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Snapshots",
		LongDesc:   "Use this command to get a list of Snapshots.",
		Example:    listSnapshotsExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunSnapshotList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace:  "snapshot",
		Resource:   "snapshot",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Snapshot",
		LongDesc:   "Use this command to get information about a specified Snapshot.\n\nRequired values to run command:\n\n* Snapshot Id",
		Example:    getSnapshotExample,
		PreCmdRun:  PreRunSnapshotId,
		CmdRun:     RunSnapshotGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv6.ArgSnapshotId, cloudapiv6.ArgIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace: "snapshot",
		Resource:  "snapshot",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Snapshot of a Volume within the Virtual Data Center",
		LongDesc: `Use this command to create a Snapshot. Creation of Snapshots is performed from the perspective of the storage Volume. The name, description and licence type of the Snapshot can be set.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Volume Id`,
		Example:    createSnapshotExample,
		PreCmdRun:  PreRunDcVolumeIds,
		CmdRun:     RunSnapshotCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Snapshot", "Name of the Snapshot")
	create.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Description of the Snapshot")
	create.AddStringFlag(cloudapiv6.ArgLicenceType, "", "LINUX", "Licence Type of the Snapshot")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"WINDOWS", "WINDOWS2016", "LINUX", "OTHER", "UNKNOWN"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgVolumeId, "", "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(cloudapiv6.ArgSecAuthProtection, "", false, "Enable secure authentication protection")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Snapshot creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace: "snapshot",
		Resource:  "snapshot",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Snapshot",
		LongDesc: `Use this command to update a specified Snapshot.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Snapshot Id`,
		Example:    updateSnapshotExample,
		PreCmdRun:  PreRunSnapshotId,
		CmdRun:     RunSnapshotUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Snapshot")
	update.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Description of the Snapshot")
	update.AddStringFlag(cloudapiv6.ArgLicenceType, "", "", "Licence Type of the Snapshot")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"WINDOWS", "WINDOWS2016", "LINUX", "OTHER", "UNKNOWN"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgSnapshotId, cloudapiv6.ArgIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", false, "This volume is capable of CPU hot plug (no reboot required)")
	update.AddBoolFlag(cloudapiv6.ArgCpuHotUnplug, "", false, "This volume is capable of CPU hot unplug (no reboot required)")
	update.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", false, "This volume is capable of memory hot plug (no reboot required)")
	update.AddBoolFlag(cloudapiv6.ArgRamHotUnplug, "", false, "This volume is capable of memory hot unplug (no reboot required)")
	update.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", false, "This volume is capable of NIC hot plug (no reboot required)")
	update.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "This volume is capable of NIC hot unplug (no reboot required)")
	update.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", false, "This volume is capable of VirtIO drive hot plug (no reboot required)")
	update.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "This volume is capable of VirtIO drive hot unplug (no reboot required)")
	update.AddBoolFlag(cloudapiv6.ArgDiscScsiHotPlug, "", false, "This volume is capable of SCSI drive hot plug (no reboot required)")
	update.AddBoolFlag(cloudapiv6.ArgDiscScsiHotUnplug, "", false, "This volume is capable of SCSI drive hot unplug (no reboot required)")
	update.AddBoolFlag(cloudapiv6.ArgSecAuthProtection, "", false, "Enable secure authentication protection")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Snapshot creation to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot creation [seconds]")

	/*
		Restore Command
	*/
	restore := core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace:  "snapshot",
		Resource:   "snapshot",
		Verb:       "restore",
		Aliases:    []string{"r"},
		ShortDesc:  "Restore a Snapshot onto a Volume",
		LongDesc:   "Use this command to restore a Snapshot onto a Volume. A Snapshot is created as just another image that can be used to create new Volumes or to restore an existing Volume.\n\nRequired values to run command:\n\n* Datacenter Id\n* Volume Id\n* Snapshot Id",
		Example:    restoreSnapshotExample,
		PreCmdRun:  PreRunSnapshotIdDcIdVolumeId,
		CmdRun:     RunSnapshotRestore,
		InitClient: true,
	})
	restore.AddStringFlag(cloudapiv6.ArgSnapshotId, cloudapiv6.ArgIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = restore.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	restore.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = restore.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	restore.AddStringFlag(cloudapiv6.ArgVolumeId, "", "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = restore.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(restore.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	restore.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Snapshot restore to be executed")
	restore.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot restore [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace:  "snapshot",
		Resource:   "snapshot",
		Verb:       "delete",
		Aliases:    []string{"d"},
		ShortDesc:  "Delete a Snapshot",
		LongDesc:   "Use this command to delete the specified Snapshot.\n\nRequired values to run command:\n\n* Snapshot Id",
		Example:    deleteSnapshotExample,
		PreCmdRun:  PreRunSnapshotDelete,
		CmdRun:     RunSnapshotDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgSnapshotId, cloudapiv6.ArgIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Snapshot deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all the Snapshots.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot deletion [seconds]")

	return snapshotCmd
}

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
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(nil, c, getSnapshots(ss)))
}

func RunSnapshotGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Snapshot with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))
	s, resp, err := c.CloudApiV6Services.Snapshots().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(nil, c, getSnapshot(s)))
}

func RunSnapshotCreate(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))
	licenseType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType))
	secAuthProtection := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgSecAuthProtection))
	c.Printer.Verbose("Properties set for creating the Snapshot: DatacenterId: %v, VolumeId: %v, Name: %v, Description: %v, LicenseType: %v, SecAuthProtection: %v",
		dcId, volumeId, name, description, licenseType, secAuthProtection)

	s, resp, err := c.CloudApiV6Services.Snapshots().Create(dcId, volumeId, name, description, licenseType, secAuthProtection)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, getSnapshot(s)))
}

func RunSnapshotUpdate(c *core.CommandConfig) error {
	c.Printer.Verbose("Updating Snapshot with id: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))
	s, resp, err := c.CloudApiV6Services.Snapshots().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)), getSnapshotPropertiesSet(c))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, getSnapshot(s)))
}

func RunSnapshotRestore(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "restore snapshot"); err != nil {
		return err
	}
	c.Printer.Verbose("Snapshot with id: %v is restoring...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))
	resp, err := c.CloudApiV6Services.Snapshots().Restore(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, nil))
}

func RunSnapshotDelete(c *core.CommandConfig) error {
	var resp *resources.Response
	var err error
	var snapshots resources.Snapshots
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll))
	if allFlag {
		_ = c.Printer.Print("Snapshots to be deleted:")
		snapshots, resp, err = c.CloudApiV6Services.Snapshots().List()
		if err != nil {
			return err
		}
		if snapshotsItems, ok := snapshots.GetItemsOk(); ok && snapshotsItems != nil {
			for _, snapshot := range *snapshotsItems {
				if id, ok := snapshot.GetIdOk(); ok && id != nil {
					_ = c.Printer.Print("Snapshot Id: " + *id)
				}
				if properties, ok := snapshot.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						_ = c.Printer.Print(" Snapshot Name: " + *name)
					}
				}
			}

			if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Snapshots"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the Snapshots...")

			for _, snapshot := range *snapshotsItems {
				if id, ok := snapshot.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Deleting Snapshot with id: %v...", *id)
					resp, err = c.CloudApiV6Services.Snapshots().Delete(*id)
					if resp != nil {
						c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
					}
					if err != nil {
						return err
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						return err
					}
				}
			}
		}
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete snapshot"); err != nil {
			return err
		}
		c.Printer.Verbose("Snapshot with id: %v is deleting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))
		resp, err := c.CloudApiV6Services.Snapshots().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))
		if resp != nil {
			c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, nil))
}

func getSnapshotPropertiesSet(c *core.CommandConfig) resources.SnapshotProperties {
	input := resources.SnapshotProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDescription)) {
		description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))
		input.SetDescription(description)
		c.Printer.Verbose("Property Description set: %v", description)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType)) {
		licenceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType))
		input.SetLicenceType(licenceType)
		c.Printer.Verbose("Property LicenceType set: %v", licenceType)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug))
		input.SetCpuHotPlug(cpuHotPlug)
		c.Printer.Verbose("Property CpuHotPlug set: %v", cpuHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotUnplug)) {
		cpuHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotUnplug))
		input.SetCpuHotUnplug(cpuHotUnplug)
		c.Printer.Verbose("Property CpuHotUnplug set: %v", cpuHotUnplug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug))
		input.SetRamHotPlug(ramHotPlug)
		c.Printer.Verbose("Property RamHotPlug set: %v", ramHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotUnplug)) {
		ramHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotUnplug))
		input.SetRamHotUnplug(ramHotUnplug)
		c.Printer.Verbose("Property RamHotUnplug set: %v", ramHotUnplug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug))
		input.SetNicHotPlug(nicHotPlug)
		c.Printer.Verbose("Property NicHotPlug set: %v", nicHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug))
		input.SetNicHotUnplug(nicHotUnplug)
		c.Printer.Verbose("Property nicHotUnplug set: %v", nicHotUnplug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug))
		input.SetDiscVirtioHotPlug(discVirtioHotPlug)
		c.Printer.Verbose("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug))
		input.SetDiscVirtioHotUnplug(discVirtioHotUnplug)
		c.Printer.Verbose("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscScsiHotPlug)) {
		discScsiHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscScsiHotPlug))
		input.SetDiscScsiHotPlug(discScsiHotPlug)
		c.Printer.Verbose("Property DiscScsiHotPlug set: %v", discScsiHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscScsiHotUnplug)) {
		discScsiHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscScsiHotUnplug))
		input.SetDiscScsiHotUnplug(discScsiHotUnplug)
		c.Printer.Verbose("Property DiscScsiHotUnplug set: %v", discScsiHotUnplug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSecAuthProtection)) {
		secAuthProtection := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgSecAuthProtection))
		input.SetSecAuthProtection(secAuthProtection)
		c.Printer.Verbose("Property SecAuthProtection set: %v", secAuthProtection)
	}
	return input
}

// Output Printing

var defaultSnapshotCols = []string{"SnapshotId", "Name", "LicenceType", "Size", "State"}

type SnapshotPrint struct {
	SnapshotId  string  `json:"SnapshotId,omitempty"`
	Name        string  `json:"Name,omitempty"`
	LicenceType string  `json:"LicenceType,omitempty"`
	Size        float32 `json:"Size,omitempty"`
	State       string  `json:"State,omitempty"`
}

func getSnapshotPrint(resp *resources.Response, c *core.CommandConfig, s []resources.Snapshot) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getSnapshotsKVMaps(s)
			r.Columns = getSnapshotCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getSnapshotCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultSnapshotCols
	}

	columnsMap := map[string]string{
		"SnapshotId":  "SnapshotId",
		"Name":        "Name",
		"LicenceType": "LicenceType",
		"Size":        "Size",
		"State":       "State",
	}
	var datacenterCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			datacenterCols = append(datacenterCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return datacenterCols
}

func getSnapshots(snapshots resources.Snapshots) []resources.Snapshot {
	ss := make([]resources.Snapshot, 0)
	if items, ok := snapshots.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ss = append(ss, resources.Snapshot{Snapshot: s})
		}
	}
	return ss
}

func getSnapshot(s *resources.Snapshot) []resources.Snapshot {
	ss := make([]resources.Snapshot, 0)
	if s != nil {
		ss = append(ss, resources.Snapshot{Snapshot: s.Snapshot})
	}
	return ss
}

func getSnapshotsKVMaps(ss []resources.Snapshot) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		o := getSnapshotKVMap(s)
		out = append(out, o)
	}
	return out
}

func getSnapshotKVMap(s resources.Snapshot) map[string]interface{} {
	var ssPrint SnapshotPrint
	if ssId, ok := s.GetIdOk(); ok && ssId != nil {
		ssPrint.SnapshotId = *ssId
	}
	if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
		if name, ok := properties.GetNameOk(); ok && name != nil {
			ssPrint.Name = *name
		}
		if licenceType, ok := properties.GetLicenceTypeOk(); ok && licenceType != nil {
			ssPrint.LicenceType = *licenceType
		}
		if size, ok := properties.GetSizeOk(); ok && size != nil {
			ssPrint.Size = *size
		}
	}
	if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			ssPrint.State = *state
		}
	}
	return structs.Map(ssPrint)
}
