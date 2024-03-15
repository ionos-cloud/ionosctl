package commands

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultLabelCols = []string{"URN", "Key", "Value", "ResourceType", "ResourceId"}
)

func LabelCmd() *core.Command {
	ctx := context.TODO()
	labelCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "label",
			Short:            "Label Operations",
			Long:             "The sub-commands of `ionosctl label` allow you to get, list, add, remove Labels from a Resource.",
			TraverseChildren: true,
		},
	}
	globalFlags := labelCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultLabelCols, tabheaders.ColsMessage(defaultLabelCols))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultLabelCols, cobra.ShellCompDirectiveNoFileComp
	})

	var (
		allowedValues = []string{cloudapiv6.DatacenterResource, cloudapiv6.VolumeResource, cloudapiv6.ServerResource, cloudapiv6.SnapshotResource, cloudapiv6.IpBlockResource}
	)

	/*
		List Command
	*/
	list := core.NewCommand(ctx, labelCmd, core.CommandBuilder{
		Namespace:  "label",
		Resource:   "label",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Labels from Resources",
		LongDesc:   "Use this command to list all Labels from all Resources under your account. If you want to list all Labels from a specific Resource, use `--resource-type` option together with the Resource Id: `--datacenter-id`, `--server-id`, `--volume-id`.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.LabelsFiltersUsage(),
		Example:    listLabelsExample,
		PreCmdRun:  PreRunLabelList,
		CmdRun:     RunLabelList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", cloudapiv6.VolumeId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", cloudapiv6.IpBlockId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgSnapshotId, "", "", cloudapiv6.SnapshotId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddSetFlag(cloudapiv6.ArgResourceType, "", "", allowedValues, "Type of resource to list labels from", core.RequiredFlagOption())
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LabelsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LabelsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	/*
		Get Command
	*/
	get := core.NewCommand(ctx, labelCmd, core.CommandBuilder{
		Namespace:  "label",
		Resource:   "label",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Label",
		LongDesc:   "Use this command to get information about a specified Label from a specified Resource.\n\nRequired values to run command:\n\n* Resource Type\n* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id\n* Label Key",
		Example:    getLabelExample,
		PreCmdRun:  PreRunResourceTypeLabelKey,
		CmdRun:     RunLabelGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv6.ArgLabelKey, "", "", cloudapiv6.LabelKey, core.RequiredFlagOption())
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", cloudapiv6.VolumeId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", cloudapiv6.IpBlockId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgSnapshotId, "", "", cloudapiv6.SnapshotId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddSetFlag(cloudapiv6.ArgResourceType, "", "", allowedValues, "Type of resource to get labels from", core.RequiredFlagOption())
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		Get By Urn Command
	*/
	getByUrn := core.NewCommand(ctx, labelCmd, core.CommandBuilder{
		Namespace:  "label",
		Resource:   "label",
		Verb:       "get-by-urn",
		ShortDesc:  "Get a Label using URN",
		LongDesc:   "Use this command to get information about a specified Label using its URN. A URN is used for uniqueness of a Label and composed using `urn:label:<resource_type>:<resource_uuid>:<key>`.\n\nRequired values to run command:\n\n* Label URN",
		Example:    getLabelByUrnExample,
		PreCmdRun:  PreRunLabelUrn,
		CmdRun:     RunLabelGetByUrn,
		InitClient: true,
	})
	getByUrn.AddStringFlag(cloudapiv6.ArgLabelUrn, "", "", "URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>]", core.RequiredFlagOption())
	getByUrn.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		Add Command
	*/
	addLabel := core.NewCommand(ctx, labelCmd, core.CommandBuilder{
		Namespace:  "label",
		Resource:   "label",
		Verb:       "add",
		Aliases:    []string{"a"},
		ShortDesc:  "Add a Label to a Resource",
		LongDesc:   "Use this command to add a Label to a specific Resource.\n\nRequired values to run command:\n\n* Resource Type\n* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id\n* Label Key\n* Label Value",
		Example:    addLabelExample,
		PreCmdRun:  PreRunResourceTypeLabelKeyValue,
		CmdRun:     RunLabelAdd,
		InitClient: true,
	})
	addLabel.AddStringFlag(cloudapiv6.ArgLabelKey, "", "", cloudapiv6.LabelKey, core.RequiredFlagOption())
	addLabel.AddStringFlag(cloudapiv6.ArgLabelValue, "", "", cloudapiv6.LabelValue, core.RequiredFlagOption())
	addLabel.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(addLabel.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", cloudapiv6.VolumeId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(addLabel.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", cloudapiv6.IpBlockId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddUUIDFlag(cloudapiv6.ArgSnapshotId, "", "", cloudapiv6.SnapshotId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddSetFlag(cloudapiv6.ArgResourceType, "", "", allowedValues, "Type of resource to add labels to", core.RequiredFlagOption())
	addLabel.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	/*
		Remove Command
	*/
	removeLabel := core.NewCommand(ctx, labelCmd, core.CommandBuilder{
		Namespace:  "label",
		Resource:   "label",
		Verb:       "remove",
		Aliases:    []string{"r"},
		ShortDesc:  "Remove a Label from a Resource",
		LongDesc:   "Use this command to remove a Label from a Resource.\n\nRequired values to run command:\n\n* Resource Type\n* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id\n* Label Key",
		Example:    removeLabelExample,
		PreCmdRun:  PreRunResourceTypeLabelKeyRemove,
		CmdRun:     RunLabelRemove,
		InitClient: true,
	})
	removeLabel.AddStringFlag(cloudapiv6.ArgLabelKey, "", "", cloudapiv6.LabelKey, core.RequiredFlagOption())
	removeLabel.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(removeLabel.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", cloudapiv6.VolumeId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(removeLabel.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", cloudapiv6.IpBlockId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddUUIDFlag(cloudapiv6.ArgSnapshotId, "", "", cloudapiv6.SnapshotId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddSetFlag(cloudapiv6.ArgResourceType, "", "", allowedValues, "Type of resource to remove labels from", core.RequiredFlagOption())
	removeLabel.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all Labels")
	removeLabel.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	return labelCmd
}

// Returns []core.FlagNameSetWithPredicate to be used as params to send to core.CheckRequiredFlagsSets funcs.
// If --resource-type datacenter, --datacenter-id is also required
// If --resource-type server,	  --datacenter-id and --server-id are also required
func generateFlagSets(c *core.PreCommandConfig, extraFlags ...string) []core.FlagNameSetWithPredicate {
	funcResourceTypeSetAndMatches := func(resource interface{}) bool {
		argResourceType := core.GetFlagName(c.NS, cloudapiv6.ArgResourceType)
		return !viper.IsSet(argResourceType) || viper.GetString(argResourceType) == resource
	}

	return []core.FlagNameSetWithPredicate{
		{
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgDataCenterId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.DatacenterResource,
		}, {
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgVolumeId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.VolumeResource,
		}, {
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.ServerResource,
		}, {
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgSnapshotId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.SnapshotResource,
		}, {
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgIpBlockId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.IpBlockResource,
		},
	}
}

func PreRunResourceTypeLabelKey(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSetsIfPredicate(c.Command, c.NS, generateFlagSets(c, cloudapiv6.ArgLabelKey)...)
}

func PreRunResourceTypeLabelKeyRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSetsIfPredicate(c.Command, c.NS,
		append(
			generateFlagSets(c, cloudapiv6.ArgLabelKey),
			generateFlagSets(c, cloudapiv6.ArgAll)...,
		)...,
	)
}

func PreRunResourceTypeLabelKeyValue(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSetsIfPredicate(c.Command, c.NS, generateFlagSets(c, cloudapiv6.ArgLabelKey, cloudapiv6.ArgLabelValue)...)
}

func PreRunLabelUrn(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgLabelUrn)
}

func PreRunLabelList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		err := query.ValidateFilters(c, completer.LabelsFilters(), completer.LabelsFiltersUsage())
		if err != nil {
			return err
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResourceType)) {
		return core.CheckRequiredFlagsSetsIfPredicate(c.Command, c.NS, generateFlagSets(c)...)
	}
	return core.NoPreRun(c)
}

func RunLabelList(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	var out string
	switch viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceType)) {
	case cloudapiv6.DatacenterResource:
		return RunDataCenterLabelsList(c)
	case cloudapiv6.ServerResource:
		return RunServerLabelsList(c)
	case cloudapiv6.VolumeResource:
		return RunVolumeLabelsList(c)
	case cloudapiv6.IpBlockResource:
		return RunIpBlockLabelsList(c)
	case cloudapiv6.SnapshotResource:
		return RunSnapshotLabelsList(c)
	default:
		labelDcs, _, err := c.CloudApiV6Services.Labels().List(listQueryParams)
		if err != nil {
			return err
		}

		cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

		out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Label, labelDcs.Labels,
			tabheaders.GetHeadersAllDefault(defaultLabelCols, cols))
		if err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

		return nil
	}
}

const labelResourceWarning = "Please use `--resource-type` flag with one option: \"datacenter\", \"volume\", \"server\", \"snapshot\", \"ipblock\""

func RunLabelGet(c *core.CommandConfig) error {
	resourceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceType))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting label with label key: %v and label value: %v for %v...", labelKey, labelValue, resourceType))

	switch resourceType {
	case cloudapiv6.DatacenterResource:
		return RunDataCenterLabelGet(c)
	case cloudapiv6.ServerResource:
		return RunServerLabelGet(c)
	case cloudapiv6.VolumeResource:
		return RunVolumeLabelGet(c)
	case cloudapiv6.IpBlockResource:
		return RunIpBlockLabelGet(c)
	case cloudapiv6.SnapshotResource:
		return RunSnapshotLabelGet(c)
	default:
		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(labelResourceWarning))

		return nil
	}
}

func RunLabelGetByUrn(c *core.CommandConfig) error {
	urn := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelUrn))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting label with urn: %v", urn))

	labelDc, _, err := c.CloudApiV6Services.Labels().GetByUrn(urn)
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Label, labelDc.Label,
		tabheaders.GetHeadersAllDefault(defaultLabelCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLabelAdd(c *core.CommandConfig) error {
	switch viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceType)) {
	case cloudapiv6.DatacenterResource:
		return RunDataCenterLabelAdd(c)
	case cloudapiv6.ServerResource:
		return RunServerLabelAdd(c)
	case cloudapiv6.VolumeResource:
		return RunVolumeLabelAdd(c)
	case cloudapiv6.IpBlockResource:
		return RunIpBlockLabelAdd(c)
	case cloudapiv6.SnapshotResource:
		return RunSnapshotLabelAdd(c)
	default:
		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(labelResourceWarning))

		return nil
	}
}

func RunLabelRemove(c *core.CommandConfig) error {
	resourceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceType))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Label is removing from %v...", resourceType))

	switch resourceType {
	case cloudapiv6.DatacenterResource:
		return RunDataCenterLabelRemove(c)
	case cloudapiv6.ServerResource:
		return RunServerLabelRemove(c)
	case cloudapiv6.VolumeResource:
		return RunVolumeLabelRemove(c)
	case cloudapiv6.IpBlockResource:
		return RunIpBlockLabelRemove(c)
	case cloudapiv6.SnapshotResource:
		return RunSnapshotLabelRemove(c)
	default:
		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(labelResourceWarning))

		return nil
	}
}
