package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func snapshot() *core.Command {
	ctx := context.TODO()
	snapshotCmd := &core.Command{
		NS: "snapshot",
		Command: &cobra.Command{
			Use:              "snapshot",
			Short:            "Snapshot Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl snapshot` + "`" + ` allow you to see information, to create, update, delete Snapshots.`,
			TraverseChildren: true,
		},
	}
	globalFlags := snapshotCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultSnapshotCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(core.GetGlobalFlagName(snapshotCmd.NS, config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace:  "snapshot",
		Resource:   "snapshot",
		Verb:       "list",
		ShortDesc:  "List Snapshots",
		LongDesc:   "Use this command to get a list of Snapshots.",
		Example:    listSnapshotsExample,
		PreCmdRun:  noPreRun,
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
		ShortDesc:  "Get a Snapshot",
		LongDesc:   "Use this command to get information about a specified Snapshot.\n\nRequired values to run command:\n\n* Snapshot Id",
		Example:    getSnapshotExample,
		PreCmdRun:  PreRunSnapshotId,
		CmdRun:     RunSnapshotGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace: "snapshot",
		Resource:  "snapshot",
		Verb:      "create",
		ShortDesc: "Create a Snapshot of a Volume within the Virtual Data Center.",
		LongDesc: `Use this command to create a Snapshot. Creation of Snapshots is performed from the perspective of the storage Volume. The name, description and licence type of the Snapshot can be set.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Volume Id
* Snapshot Name
* Snapshot Licence Type`,
		Example:    createSnapshotExample,
		PreCmdRun:  PreRunSnapNameLicenceDcIdVolumeId,
		CmdRun:     RunSnapshotCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgSnapshotName, "", "", "Name of the Snapshot"+config.RequiredFlag)
	create.AddStringFlag(config.ArgSnapshotDescription, "", "", "Description of the Snapshot")
	create.AddStringFlag(config.ArgSnapshotLicenceType, "", "", "Licence Type of the Snapshot"+config.RequiredFlag)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgSnapshotLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"WINDOWS", "WINDOWS2016", "LINUX", "OTHER", "UNKNOWN"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgSnapshotSecAuthProtection, "", false, "Enable secure authentication protection")
	create.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for Snapshot creation to be executed")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace: "snapshot",
		Resource:  "snapshot",
		Verb:      "update",
		ShortDesc: "Update a Snapshot.",
		LongDesc: `Use this command to update a specified Snapshot.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Snapshot Id`,
		Example:    updateSnapshotExample,
		PreCmdRun:  PreRunSnapshotId,
		CmdRun:     RunSnapshotUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgSnapshotName, "", "", "Name of the Snapshot")
	update.AddStringFlag(config.ArgSnapshotDescription, "", "", "Description of the Snapshot")
	update.AddStringFlag(config.ArgSnapshotLicenceType, "", "", "Licence Type of the Snapshot")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgSnapshotLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"WINDOWS", "WINDOWS2016", "LINUX", "OTHER", "UNKNOWN"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgSnapshotCpuHotPlug, "", false, "This volume is capable of CPU hot plug (no reboot required)")
	update.AddBoolFlag(config.ArgSnapshotCpuHotUnplug, "", false, "This volume is capable of CPU hot unplug (no reboot required)")
	update.AddBoolFlag(config.ArgSnapshotRamHotPlug, "", false, "This volume is capable of memory hot plug (no reboot required)")
	update.AddBoolFlag(config.ArgSnapshotRamHotUnplug, "", false, "This volume is capable of memory hot unplug (no reboot required)")
	update.AddBoolFlag(config.ArgSnapshotNicHotPlug, "", false, "This volume is capable of NIC hot plug (no reboot required)")
	update.AddBoolFlag(config.ArgSnapshotNicHotUnplug, "", false, "This volume is capable of NIC hot unplug (no reboot required)")
	update.AddBoolFlag(config.ArgSnapshotDiscVirtioHotPlug, "", false, "This volume is capable of VirtIO drive hot plug (no reboot required)")
	update.AddBoolFlag(config.ArgSnapshotDiscVirtioHotUnplug, "", false, "This volume is capable of VirtIO drive hot unplug (no reboot required)")
	update.AddBoolFlag(config.ArgSnapshotDiscScsiHotPlug, "", false, "This volume is capable of SCSI drive hot plug (no reboot required)")
	update.AddBoolFlag(config.ArgSnapshotDiscScsiHotUnplug, "", false, "This volume is capable of SCSI drive hot unplug (no reboot required)")
	update.AddBoolFlag(config.ArgSnapshotSecAuthProtection, "", false, "Enable secure authentication protection")
	update.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for Snapshot creation to be executed")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot creation [seconds]")

	/*
		Restore Command
	*/
	restore := core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace:  "snapshot",
		Resource:   "snapshot",
		Verb:       "restore",
		ShortDesc:  "Restore a Snapshot onto a Volume",
		LongDesc:   "Use this command to restore a Snapshot onto a Volume. A Snapshot is created as just another image that can be used to create new Volumes or to restore an existing Volume.\n\nRequired values to run command:\n\n* Datacenter Id\n* Volume Id\n* Snapshot Id",
		Example:    restoreSnapshotExample,
		PreCmdRun:  PreRunSnapshotIdDcIdVolumeId,
		CmdRun:     RunSnapshotRestore,
		InitClient: true,
	})
	restore.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = restore.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	restore.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = restore.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	restore.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = restore.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(restore.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	restore.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for Snapshot restore to be executed")
	restore.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot restore [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace:  "snapshot",
		Resource:   "snapshot",
		Verb:       "delete",
		ShortDesc:  "Delete a Snapshot",
		LongDesc:   "Use this command to delete the specified Snapshot.\n\nRequired values to run command:\n\n* Snapshot Id",
		Example:    deleteSnapshotExample,
		PreCmdRun:  PreRunSnapshotId,
		CmdRun:     RunSnapshotDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for Snapshot deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot deletion [seconds]")

	return snapshotCmd
}

func PreRunSnapshotId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgSnapshotId)
}

func PreRunSnapNameLicenceDcIdVolumeId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgVolumeId, config.ArgSnapshotName, config.ArgSnapshotLicenceType)
}

func PreRunSnapshotIdDcIdVolumeId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgVolumeId, config.ArgSnapshotId)
}

func RunSnapshotList(c *core.CommandConfig) error {
	ss, _, err := c.Snapshots().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(nil, c, getSnapshots(ss)))
}

func RunSnapshotGet(c *core.CommandConfig) error {
	s, _, err := c.Snapshots().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(nil, c, getSnapshot(s)))
}

func RunSnapshotCreate(c *core.CommandConfig) error {
	s, resp, err := c.Snapshots().Create(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotName)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotDescription)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotLicenceType)),
		viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotSecAuthProtection)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, getSnapshot(s)))
}

func RunSnapshotUpdate(c *core.CommandConfig) error {
	s, resp, err := c.Snapshots().Update(viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotId)), getSnapshotPropertiesSet(c))
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, getSnapshot(s)))
}

func RunSnapshotRestore(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "restore snapshot"); err != nil {
		return err
	}
	resp, err := c.Snapshots().Restore(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, nil))
}

func RunSnapshotDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete snapshot"); err != nil {
		return err
	}
	resp, err := c.Snapshots().Delete(viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, nil))
}

func getSnapshotPropertiesSet(c *core.CommandConfig) resources.SnapshotProperties {
	input := resources.SnapshotProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotDescription)) {
		input.SetDescription(viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotDescription)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotLicenceType)) {
		input.SetLicenceType(viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotLicenceType)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotCpuHotPlug)) {
		input.SetCpuHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotCpuHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotCpuHotUnplug)) {
		input.SetCpuHotUnplug(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotCpuHotUnplug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotRamHotPlug)) {
		input.SetRamHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotRamHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotRamHotUnplug)) {
		input.SetRamHotUnplug(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotRamHotUnplug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotNicHotPlug)) {
		input.SetNicHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotNicHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotNicHotUnplug)) {
		input.SetNicHotUnplug(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotNicHotUnplug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotDiscVirtioHotPlug)) {
		input.SetDiscVirtioHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotDiscVirtioHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotDiscVirtioHotUnplug)) {
		input.SetDiscVirtioHotUnplug(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotDiscVirtioHotUnplug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotDiscScsiHotPlug)) {
		input.SetDiscScsiHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotDiscScsiHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotDiscScsiHotUnplug)) {
		input.SetDiscScsiHotUnplug(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotDiscScsiHotUnplug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSnapshotSecAuthProtection)) {
		input.SetSecAuthProtection(viper.GetBool(core.GetFlagName(c.NS, config.ArgSnapshotSecAuthProtection)))
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
			r.Columns = getSnapshotCols(core.GetGlobalFlagName(c.Namespace, config.ArgCols), c.Printer.GetStderr())
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

func getSnapshotIds(outErr io.Writer) []string {
	err := config.LoadFile()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	snapshotSvc := resources.NewSnapshotService(clientSvc.Get(), context.TODO())
	snapshots, _, err := snapshotSvc.List()
	clierror.CheckError(err, outErr)
	ssIds := make([]string, 0)
	if items, ok := snapshots.Snapshots.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}
