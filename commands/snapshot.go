package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func snapshot() *builder.Command {
	ctx := context.TODO()
	snapshotCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "snapshot",
			Aliases:          []string{"sp", "snap"},
			Short:            "Snapshot Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl snapshot` + "`" + ` allows you to see information about Snapshots.`,
			TraverseChildren: true,
		},
	}
	globalFlags := snapshotCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultSnapshotCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(snapshotCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, snapshotCmd, noPreRun, RunSnapshotList, "list", "List Snapshots",
		"Use this command to get a list of Snapshots.", listSnapshotsExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, snapshotCmd, noPreRun, RunSnapshotGet, "get", "Get a Snapshot",
		"Use this command to get information about a specified Snapshot.\n\nRequired values to run command:\n- Snapshot Id",
		getSnapshotExample, true)
	get.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, snapshotCmd, PreRunSnapshotNameDcIdVolumeIdValidate, RunSnapshotCreate, "create", "Create a Snapshot of a Volume within the Virtual Data Center.",
		`Use this command to create a Snapshot in a specified Data Center. Creation of Snapshots is performed from the perspective of the storage volume. The name, description and licence type of the Snapshot can be set.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:
- Data Center Id
- Volume Id
- Snapshot Name`, createSnapshotExample, true)
	create.AddStringFlag(config.ArgSnapshotName, "", "", "Name of the Snapshot")
	create.AddStringFlag(config.ArgSnapshotDescription, "", "", "Description of the Snapshot")
	create.AddStringFlag(config.ArgSnapshotLicenceType, "", "", "Licence Type of the Snapshot")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgSnapshotLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"WINDOWS", "WINDOWS2016", "LINUX", "OTHER", "UNKNOWN"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetFlagName(snapshotCmd.Command.Name(), create.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Snapshot to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for a Snapshot to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, snapshotCmd, PreRunSnapshotIdValidate, RunSnapshotUpdate, "update", "Update a Snapshot.",
		`Use this command to update a specified Snapshot.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:
- Snapshot Id`, updateSnapshotExample, true)
	update.AddStringFlag(config.ArgSnapshotName, "", "", "Name of the Snapshot")
	update.AddStringFlag(config.ArgSnapshotDescription, "", "", "Description of the Snapshot")
	update.AddFloat32Flag(config.ArgSnapshotSize, "", config.DefaultVolumeSize, "Size of the Snapshot")
	update.AddStringFlag(config.ArgSnapshotLicenceType, "", "", "Licence Type of the Snapshot")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgSnapshotLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"WINDOWS", "WINDOWS2016", "LINUX", "OTHER", "UNKNOWN"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Snapshot to be created")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for a Snapshot to be created [seconds]")

	/*
		Restore Command
	*/
	restore := builder.NewCommand(ctx, snapshotCmd, PreRunSnapshotIdDcIdVolumeIdValidate, RunSnapshotRestore, "restore", "Restore a Snapshot onto a Volume",
		"Use this command to restore a Snapshot onto a Volume. A Snapshot is created as just another image that can be used to create new Volumes or to restore an existing Volume.\n\nRequired values to run command:\n\n* Datacenter Id\n* Volume Id\n* Snapshot Id",
		restoreSnapshotExample, true)
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
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetFlagName(snapshotCmd.Command.Name(), restore.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	restore.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Snapshot to be restored")
	restore.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for a Snapshot to be restored [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, snapshotCmd, PreRunSnapshotIdValidate, RunSnapshotDelete, "delete", "Delete a Snapshot",
		"Use this command to delete the specified Snapshot.\n\nRequired values to run command:\n\n* Snapshot Id",
		deleteSnapshotExample, true)
	deleteCmd.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Snapshot to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for a Snapshot to be deleted [seconds]")

	return snapshotCmd
}

func PreRunSnapshotIdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgSnapshotId)
}

func PreRunSnapshotNameDcIdVolumeIdValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgDataCenterId, config.ArgVolumeId, config.ArgSnapshotName)
	if err != nil {
		return err
	}
	return nil
}

func PreRunSnapshotIdDcIdVolumeIdValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgDataCenterId, config.ArgVolumeId, config.ArgSnapshotId)
	if err != nil {
		return err
	}
	return nil
}

func RunSnapshotList(c *builder.CommandConfig) error {
	ss, _, err := c.Snapshots().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(nil, c, getSnapshots(ss)))
}

func RunSnapshotGet(c *builder.CommandConfig) error {
	s, _, err := c.Snapshots().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(nil, c, getSnapshot(s)))
}

func RunSnapshotCreate(c *builder.CommandConfig) error {
	datacenterId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId))
	volumeId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId))
	name := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotName))
	description := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotDescription))
	licenceType := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotLicenceType))
	s, resp, err := c.Snapshots().Create(datacenterId, volumeId, name, description, licenceType)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, getSnapshot(s)))
}

func RunSnapshotUpdate(c *builder.CommandConfig) error {
	input := resources.SnapshotProperties{}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotName)) {
		input.SetName(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotName)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotDescription)) {
		input.SetDescription(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotDescription)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotLicenceType)) {
		input.SetLicenceType(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotLicenceType)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotSize)) {
		input.SetSize(float32(viper.GetFloat64(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotSize))))
	}
	s, resp, err := c.Snapshots().Update(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId)),
		input,
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, getSnapshot(s)))
}

func RunSnapshotRestore(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "restore snapshot")
	if err != nil {
		return err
	}
	datacenterId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId))
	volumeId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId))
	snapshotId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId))
	resp, err := c.Snapshots().Restore(datacenterId, volumeId, snapshotId)
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, nil))
}

func RunSnapshotDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete snapshot")
	if err != nil {
		return err
	}
	snapshotId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId))
	resp, err := c.Snapshots().Delete(snapshotId)
	if err != nil {
		return err
	}
	return c.Printer.Print(getSnapshotPrint(resp, c, nil))
}

var defaultSnapshotCols = []string{"SnapshotId", "Name", "LicenceType", "Size"}

type SnapshotPrint struct {
	SnapshotId  string  `json:"SnapshotId,omitempty"`
	Name        string  `json:"Name,omitempty"`
	LicenceType string  `json:"LicenceType,omitempty"`
	Size        float32 `json:"Size,omitempty"`
}

func getSnapshotPrint(resp *resources.Response, c *builder.CommandConfig, s []resources.Snapshot) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getSnapshotsKVMaps(s)
			r.Columns = getSnapshotCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
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
	if s != nil {
		return []resources.Snapshot{*s}
	}
	return nil
}

func getSnapshotsKVMaps(ss []resources.Snapshot) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		properties := s.GetProperties()
		var ssPrint SnapshotPrint
		if ssId, ok := s.GetIdOk(); ok && ssId != nil {
			ssPrint.SnapshotId = *ssId
		}
		if name, ok := properties.GetNameOk(); ok && name != nil {
			ssPrint.Name = *name
		}
		if licenceType, ok := properties.GetLicenceTypeOk(); ok && licenceType != nil {
			ssPrint.LicenceType = *licenceType
		}
		if size, ok := properties.GetSizeOk(); ok && size != nil {
			ssPrint.Size = *size
		}
		o := structs.Map(ssPrint)
		out = append(out, o)
	}
	return out
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
