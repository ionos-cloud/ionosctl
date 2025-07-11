package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultSnapshotCols = []string{"SnapshotId", "Name", "LicenceType", "Size", "State"}
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
	globalFlags.StringSliceP(constants.FlagCols, "", defaultSnapshotCols, tabheaders.ColsMessage(defaultSnapshotCols))
	_ = viper.BindPFlag(core.GetFlagName(snapshotCmd.Name(), constants.FlagCols), globalFlags.Lookup(constants.FlagCols))
	_ = snapshotCmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultSnapshotCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, snapshotCmd, core.CommandBuilder{
		Namespace:  "snapshot",
		Resource:   "snapshot",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Snapshots",
		LongDesc:   "Use this command to get a list of Snapshots.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.SnapshotsFiltersUsage(),
		Example:    listSnapshotsExample,
		PreCmdRun:  PreRunSnapshotList,
		CmdRun:     RunSnapshotList,
		InitClient: true,
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.FlagDepthDescription)
	list.AddStringFlag(cloudapiv6.FlagOrderBy, "", "", cloudapiv6.FlagOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.FlagFilters, cloudapiv6.FlagFiltersShort, []string{""}, cloudapiv6.FlagFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotsFilters(), cobra.ShellCompDirectiveNoFileComp
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
	get.AddUUIDFlag(cloudapiv6.FlagSnapshotId, cloudapiv6.FlagIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.FlagDepthDescription)

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
	create.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "Unnamed Snapshot", "Name of the Snapshot")
	create.AddStringFlag(cloudapiv6.FlagDescription, cloudapiv6.FlagDescriptionShort, "", "Description of the Snapshot")
	create.AddSetFlag(cloudapiv6.FlagLicenceType, "", "LINUX", constants.EnumLicenceType, "Licence Type of the Snapshot")
	create.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.FlagVolumeId, "", "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(create.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(cloudapiv6.FlagSecAuthProtection, "", false, "Enable secure authentication protection. E.g.: --sec-auth-protection=true, --sec-auth-protection=false")
	create.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Snapshot creation to be executed")
	create.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot creation [seconds]")
	create.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.FlagDepthDescription)

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
	update.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "", "Name of the Snapshot")
	update.AddStringFlag(cloudapiv6.FlagDescription, cloudapiv6.FlagDescriptionShort, "", "Description of the Snapshot")
	update.AddSetFlag(cloudapiv6.FlagLicenceType, "", "", constants.EnumLicenceType, "Licence Type of the Snapshot")
	update.AddUUIDFlag(cloudapiv6.FlagSnapshotId, cloudapiv6.FlagIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(cloudapiv6.FlagCpuHotPlug, "", false, "This volume is capable of CPU hot plug (no reboot required). E.g.: --cpu-hot-plug=true, --cpu-hot-plug=false")
	update.AddBoolFlag(cloudapiv6.FlagCpuHotUnplug, "", false, "This volume is capable of CPU hot unplug (no reboot required). E.g.: --cpu-hot-unplug=true, --cpu-hot-unplug=false")
	update.AddBoolFlag(cloudapiv6.FlagRamHotPlug, "", false, "This volume is capable of memory hot plug (no reboot required). E.g.: --ram-hot-plug=true, --ram-hot-plug=false")
	update.AddBoolFlag(cloudapiv6.FlagRamHotUnplug, "", false, "This volume is capable of memory hot unplug (no reboot required). E.g.: --ram-hot-unplug=true, --ram-hot-unplug=false")
	update.AddBoolFlag(cloudapiv6.FlagNicHotPlug, "", false, "This volume is capable of NIC hot plug (no reboot required). E.g.: --nic-hot-plug=true, --nic-hot-plug=false")
	update.AddBoolFlag(cloudapiv6.FlagNicHotUnplug, "", false, "This volume is capable of NIC hot unplug (no reboot required). E.g.: --nic-hot-unplug=true, --nic-hot-unplug=false")
	update.AddBoolFlag(cloudapiv6.FlagDiscVirtioHotPlug, "", false, "This volume is capable of VirtIO drive hot plug (no reboot required). E.g.: --disc-virtio-plug=true, --disc-virtio-plug=false")
	update.AddBoolFlag(cloudapiv6.FlagDiscVirtioHotUnplug, "", false, "This volume is capable of VirtIO drive hot unplug (no reboot required). E.g.: --disc-virtio-unplug=true, --disc-virtio-unplug=false")
	update.AddBoolFlag(cloudapiv6.FlagDiscScsiHotPlug, "", false, "This volume is capable of SCSI drive hot plug (no reboot required). E.g.: --disc-scsi-plug=true, --disc-scsi-plug=false")
	update.AddBoolFlag(cloudapiv6.FlagDiscScsiHotUnplug, "", false, "This volume is capable of SCSI drive hot unplug (no reboot required). E.g.: --disc-scsi-unplug=true, --disc-scsi-unplug=false")
	update.AddBoolFlag(cloudapiv6.FlagSecAuthProtection, "", false, "Enable secure authentication protection. E.g.: --sec-auth-protection=true, --sec-auth-protection=false")
	update.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Snapshot creation to be executed")
	update.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot creation [seconds]")
	update.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.FlagDepthDescription)

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
	restore.AddUUIDFlag(cloudapiv6.FlagSnapshotId, cloudapiv6.FlagIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = restore.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	restore.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = restore.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	restore.AddUUIDFlag(cloudapiv6.FlagVolumeId, "", "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = restore.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(restore.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	restore.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Snapshot restore to be executed")
	restore.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot restore [seconds]")
	restore.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.FlagDepthDescription)

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
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagSnapshotId, cloudapiv6.FlagIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Snapshot deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, "Delete all the Snapshots.")
	deleteCmd.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.FlagDepthDescription)

	return core.WithConfigOverride(snapshotCmd, "compute", "")
}

func PreRunSnapshotList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagFilters)) {
		return query.ValidateFilters(c, completer.SnapshotsFilters(), completer.SnapshotsFiltersUsage())
	}
	return nil
}

func PreRunSnapshotId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagSnapshotId)
}

func PreRunSnapshotDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagSnapshotId},
		[]string{cloudapiv6.FlagAll},
	)
}

func PreRunSnapshotIdDcIdVolumeId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagVolumeId, cloudapiv6.FlagSnapshotId)
}

func RunSnapshotList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	ss, resp, err := c.CloudApiV6Services.Snapshots().List(listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Snapshot, ss.Snapshots,
		tabheaders.GetHeadersAllDefault(defaultSnapshotCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunSnapshotGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Snapshot with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSnapshotId))))

	s, resp, err := c.CloudApiV6Services.Snapshots().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSnapshotId)), queryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Snapshot, s.Snapshot,
		tabheaders.GetHeadersAllDefault(defaultSnapshotCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunSnapshotCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagVolumeId))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))
	description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDescription))
	licenseType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagLicenceType))
	secAuthProtection := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagSecAuthProtection))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Properties set for creating the Snapshot: DatacenterId: %v, VolumeId: %v, Name: %v, Description: %v, "+
			"LicenseType: %v, SecAuthProtection: %v", dcId, volumeId, name, description, licenseType, secAuthProtection))

	s, resp, err := c.CloudApiV6Services.Snapshots().Create(dcId, volumeId, name, description, licenseType, secAuthProtection, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Snapshot, s.Snapshot,
		tabheaders.GetHeadersAllDefault(defaultSnapshotCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunSnapshotUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Updating Snapshot with id: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSnapshotId))))

	s, resp, err := c.CloudApiV6Services.Snapshots().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSnapshotId)),
		getSnapshotPropertiesSet(c), queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Snapshot, s.Snapshot,
		tabheaders.GetHeadersAllDefault(defaultSnapshotCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunSnapshotRestore(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "restore snapshot", viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Snapshot with id: %v is restoring...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSnapshotId))))

	resp, err := c.CloudApiV6Services.Snapshots().Restore(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagVolumeId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSnapshotId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Snapshot successfully restored"))
	return nil
}

func RunSnapshotDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSnapshotId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		if err := DeleteAllSnapshots(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete snapshot", viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Snapshot with id: %v...", snapshotId))

	resp, err := c.CloudApiV6Services.Snapshots().Delete(snapshotId, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Snapshot successfully deleted"))
	return nil
}

func getSnapshotPropertiesSet(c *core.CommandConfig) resources.SnapshotProperties {
	input := resources.SnapshotProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))
		input.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagDescription)) {
		description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDescription))
		input.SetDescription(description)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Description set: %v", description))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagLicenceType)) {
		licenceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagLicenceType))
		input.SetLicenceType(licenceType)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property LicenceType set: %v", licenceType))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagCpuHotPlug))
		input.SetCpuHotPlug(cpuHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property CpuHotPlug set: %v", cpuHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagCpuHotUnplug)) {
		cpuHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagCpuHotUnplug))
		input.SetCpuHotUnplug(cpuHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property CpuHotUnplug set: %v", cpuHotUnplug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagRamHotPlug))
		input.SetRamHotPlug(ramHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property RamHotPlug set: %v", ramHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagRamHotUnplug)) {
		ramHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagRamHotUnplug))
		input.SetRamHotUnplug(ramHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property RamHotUnplug set: %v", ramHotUnplug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagNicHotPlug))
		input.SetNicHotPlug(nicHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property NicHotPlug set: %v", nicHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagNicHotUnplug))
		input.SetNicHotUnplug(nicHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property nicHotUnplug set: %v", nicHotUnplug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagDiscVirtioHotPlug))
		input.SetDiscVirtioHotPlug(discVirtioHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagDiscVirtioHotUnplug))
		input.SetDiscVirtioHotUnplug(discVirtioHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagDiscScsiHotPlug)) {
		discScsiHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagDiscScsiHotPlug))
		input.SetDiscScsiHotPlug(discScsiHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property DiscScsiHotPlug set: %v", discScsiHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagDiscScsiHotUnplug)) {
		discScsiHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagDiscScsiHotUnplug))
		input.SetDiscScsiHotUnplug(discScsiHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property DiscScsiHotUnplug set: %v", discScsiHotUnplug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagSecAuthProtection)) {
		secAuthProtection := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagSecAuthProtection))
		input.SetSecAuthProtection(secAuthProtection)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property SecAuthProtection set: %v", secAuthProtection))
	}

	return input
}

func DeleteAllSnapshots(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Snapshots..."))

	snapshots, resp, err := c.CloudApiV6Services.Snapshots().List(cloudapiv6.ParentResourceListQueryParams)
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

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Snapshot with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.FlagForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Snapshots().Delete(*id, queryParams)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
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
