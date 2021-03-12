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
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func snapshot() *builder.Command {
	snapshotCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "snapshot",
			Short:            "Snapshot Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl snapshot` + "`" + ` allows you to see information about snapshots available.`,
			TraverseChildren: true,
		},
	}
	globalFlags := snapshotCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultSnapshotCols, "Columns to be printed in the standard output")
	viper.BindPFlag(builder.GetGlobalFlagName(snapshotCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(context.TODO(), snapshotCmd, noPreRun, RunSnapshotList, "list", "List Snapshots",
		"Use this command to get a list of available snapshots to create objects on.", "", true)

	/*
		Get Command
	*/
	get := builder.NewCommand(context.TODO(), snapshotCmd, noPreRun, RunSnapshotGet, "get", "Get a Snapshot",
		"Use this command to get information about a specified Snapshot.\n\nRequired values to run command:\n- Snapshot Id",
		"", true)
	get.AddStringFlag(config.ArgSnapshotId, "", "", "The unique Snapshot Id [Required flag]")
	get.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(context.TODO(), snapshotCmd, noPreRun, RunSnapshotCreate, "create", "Create a Snapshot",
		`Use this command to create a Server in a specified Data Center. The name, cores, ram, cpu-family and availability zone options can be set.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:
- Data Center Id`, "", true)
	create.AddStringFlag(config.ArgSnapshotName, "", "", "Name of the Server")
	create.AddStringFlag(config.ArgSnapshotDescription, "", "", "CPU Family for the Server")
	create.AddStringFlag(config.ArgSnapshotLicenceType, "", "", "CPU Family for the Server")
	create.AddStringFlag(config.ArgSnapshotId, "", "", "The unique Snapshot Id [Required flag]")
	create.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", "The unique Snapshot Id [Required flag]")
	create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgVolumeId, "", "", "The unique Snapshot Id [Required flag]")
	create.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, snapshotCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Server to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option [seconds]")

	return snapshotCmd
}

func RunSnapshotList(c *builder.CommandConfig) error {
	snapshots, _, err := c.Snapshots().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: snapshots,
		KeyValue:   getSnapshotsKVMaps(getSnapshots(snapshots)),
		Columns:    getSnapshotCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunSnapshotGet(c *builder.CommandConfig) error {
	snapshot, _, err := c.Snapshots().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: snapshot,
		KeyValue:   getSnapshotsKVMaps([]resources.Snapshot{*snapshot}),
		Columns:    getSnapshotCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunSnapshotCreate(c *builder.CommandConfig) error {
	datacenterId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId))
	volumeId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId))
	name := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotName))
	description := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotDescription))
	licencetype := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotLicenceType))
	snapshot, resp, err := c.Snapshots().Create(datacenterId, volumeId, name, description, licencetype)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "snapshot",
		Verb:        "create",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
		OutputJSON:  snapshot,
		KeyValue:    getSnapshotsKVMaps([]resources.Snapshot{*snapshot}),
		Columns:     getSnapshotCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunSnapshotRestore(c *builder.CommandConfig) error {
	datacenterId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId))
	volumeId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId))
	snapshotId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId))
	resp, err := c.Snapshots().Restore(datacenterId, volumeId, snapshotId)
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "snapshot",
		Verb:        "restore",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunSnapshotDelete(c *builder.CommandConfig) error {
	snapshotId := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId))
	resp, err := c.Snapshots().Delete(snapshotId)
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "snapshot",
		Verb:        "delete",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

var defaultSnapshotCols = []string{"SnapshotId", "Name", "LicenceType", "Size"}

type SnapshotPrint struct {
	SnapshotId  string  `json:"SnapshotId,omitempty"`
	Name        string  `json:"Name,omitempty"`
	LicenceType string  `json:"LicenceType,omitempty"`
	Size        float32 `json:"Size,omitempty"`
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

func getSnapshots(datacenters resources.Snapshots) []resources.Snapshot {
	dc := make([]resources.Snapshot, 0)
	for _, d := range *datacenters.Items {
		dc = append(dc, resources.Snapshot{d})
	}
	return dc
}

func getSnapshotsKVMaps(dcs []resources.Snapshot) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(dcs))
	for _, dc := range dcs {
		properties := dc.GetProperties()
		var dcPrint SnapshotPrint
		if dcid, ok := dc.GetIdOk(); ok && dcid != nil {
			dcPrint.SnapshotId = *dcid
		}
		if name, ok := properties.GetNameOk(); ok && name != nil {
			dcPrint.Name = *name
		}
		if licenceType, ok := properties.GetLicenceTypeOk(); ok && licenceType != nil {
			dcPrint.LicenceType = *licenceType
		}
		if size, ok := properties.GetSizeOk(); ok && size != nil {
			dcPrint.Size = *size
		}
		o := structs.Map(dcPrint)
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
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)

	snapshotSvc := resources.NewSnapshotService(clientSvc.Get(), context.TODO())
	snapshots, _, err := snapshotSvc.List()
	clierror.CheckError(err, outErr)

	lcIds := make([]string, 0)
	if snapshots.Snapshots.Items != nil {
		for _, d := range *snapshots.Snapshots.Items {
			lcIds = append(lcIds, *d.GetId())
		}
	} else {
		return nil
	}
	return lcIds
}
