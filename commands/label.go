package commands

import (
	"context"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	multierror "go.uber.org/multierr"
)

func label() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultLabelResourceCols, utils.ColsMessage(defaultLabelResourceCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(labelCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultLabelResourceCols, cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgDataCenterId, "", "", "The unique Data Center Id")
	_ = viper.BindPFlag(core.GetGlobalFlagName(labelCmd.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgServerId, "", "", "The unique Server Id")
	_ = viper.BindPFlag(core.GetGlobalFlagName(labelCmd.Name(), config.ArgServerId), globalFlags.Lookup(config.ArgServerId))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(labelCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgVolumeId, "", "", "The unique Volume Id")
	_ = viper.BindPFlag(core.GetGlobalFlagName(labelCmd.Name(), config.ArgVolumeId), globalFlags.Lookup(config.ArgVolumeId))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(labelCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgIpBlockId, "", "", "The unique IpBlock Id")
	_ = viper.BindPFlag(core.GetGlobalFlagName(labelCmd.Name(), config.ArgIpBlockId), globalFlags.Lookup(config.ArgIpBlockId))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgSnapshotId, "", "", "The unique Snapshot Id")
	_ = viper.BindPFlag(core.GetGlobalFlagName(labelCmd.Name(), config.ArgSnapshotId), globalFlags.Lookup(config.ArgSnapshotId))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgResourceType, "", "", "Resource Type")
	_ = viper.BindPFlag(core.GetGlobalFlagName(labelCmd.Name(), config.ArgResourceType), globalFlags.Lookup(config.ArgResourceType))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"datacenter", "volume", "server", "snapshot", "ipblock"}, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, labelCmd, core.CommandBuilder{
		Namespace:  "label",
		Resource:   "label",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Labels from Resources",
		LongDesc:   "Use this command to list all Labels from all Resources under your account. If you want to list all Labels from a specific Resource, use `--resource-type` option together with the Resource Id: `--datacenter-id`, `--server-id`, `--volume-id`.",
		Example:    listLabelsExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunLabelList,
		InitClient: true,
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
		PreCmdRun:  PreRunGlobalResourceTypeLabelKey,
		CmdRun:     RunLabelGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)

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
	getByUrn.AddStringFlag(config.ArgLabelUrn, "", "", "URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>] (required)")

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
		PreCmdRun:  PreRunGlobalResourceTypeLabelKeyValue,
		CmdRun:     RunLabelAdd,
		InitClient: true,
	})
	addLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	addLabel.AddStringFlag(config.ArgLabelValue, "", "", config.RequiredFlagLabelValue)

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
		PreCmdRun:  PreRunGlobalResourceTypeLabelKey,
		CmdRun:     RunLabelRemove,
		InitClient: true,
	})
	removeLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)

	return labelCmd
}

func PreRunGlobalResourceTypeLabelKey(c *core.PreCommandConfig) error {
	var result error
	if err := core.CheckRequiredGlobalFlags(c.Resource, config.ArgResourceType); err != nil {
		result = multierror.Append(result, err)
	}
	if err := checkResourceIds(c); err != nil {
		result = multierror.Append(result, err)
	}
	if err := core.CheckRequiredFlags(c.NS, config.ArgLabelKey); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func PreRunGlobalResourceTypeLabelKeyValue(c *core.PreCommandConfig) error {
	var result error
	if err := core.CheckRequiredGlobalFlags(c.Resource, config.ArgResourceType); err != nil {
		result = multierror.Append(result, err)
	}
	if err := checkResourceIds(c); err != nil {
		result = multierror.Append(result, err)
	}
	if err := core.CheckRequiredFlags(c.NS, config.ArgLabelKey, config.ArgLabelValue); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func PreRunLabelUrn(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgLabelUrn)
}

const labelResourceWarning = "Please use `--resource-type` flag with one option: \"datacenter\", \"volume\", \"server\", \"snapshot\", \"ipblock\""

func RunLabelList(c *core.CommandConfig) error {
	switch viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgResourceType)) {
	case config.DatacenterResource:
		return RunDataCenterLabelsList(c)
	case config.ServerResource:
		return RunServerLabelsList(c)
	case config.VolumeResource:
		return RunVolumeLabelsList(c)
	case config.IpBlockResource:
		return RunIpBlockLabelsList(c)
	case config.SnapshotResource:
		return RunSnapshotLabelsList(c)
	default:
		labelDcs, _, err := c.Labels().List()
		if err != nil {
			return err
		}
		return c.Printer.Print(getLabelPrint(c, getLabels(labelDcs)))
	}
}

func RunLabelGet(c *core.CommandConfig) error {
	c.Printer.Info("Getting label...")
	switch viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgResourceType)) {
	case config.DatacenterResource:
		return RunDataCenterLabelGet(c)
	case config.ServerResource:
		return RunServerLabelGet(c)
	case config.VolumeResource:
		return RunVolumeLabelGet(c)
	case config.IpBlockResource:
		return RunIpBlockLabelGet(c)
	case config.SnapshotResource:
		return RunSnapshotLabelGet(c)
	default:
		return c.Printer.Print(labelResourceWarning)
	}
}

func RunLabelGetByUrn(c *core.CommandConfig) error {
	c.Printer.Info("Getting label by urn...")
	labelDc, _, err := c.Labels().GetByUrn(viper.GetString(core.GetFlagName(c.NS, config.ArgLabelUrn)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelPrint(c, getLabel(labelDc)))
}

func RunLabelAdd(c *core.CommandConfig) error {
	switch viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgResourceType)) {
	case config.DatacenterResource:
		return RunDataCenterLabelAdd(c)
	case config.ServerResource:
		return RunServerLabelAdd(c)
	case config.VolumeResource:
		return RunVolumeLabelAdd(c)
	case config.IpBlockResource:
		return RunIpBlockLabelAdd(c)
	case config.SnapshotResource:
		return RunSnapshotLabelAdd(c)
	default:
		return c.Printer.Print(labelResourceWarning)
	}
}

func RunLabelRemove(c *core.CommandConfig) error {
	c.Printer.Info("Label is removing...")
	switch viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgResourceType)) {
	case config.DatacenterResource:
		return RunDataCenterLabelRemove(c)
	case config.ServerResource:
		return RunServerLabelRemove(c)
	case config.VolumeResource:
		return RunVolumeLabelRemove(c)
	case config.IpBlockResource:
		return RunIpBlockLabelRemove(c)
	case config.SnapshotResource:
		return RunSnapshotLabelRemove(c)
	default:
		return c.Printer.Print(labelResourceWarning)
	}
}

func checkResourceIds(c *core.PreCommandConfig) error {
	switch viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgResourceType)) {
	case config.DatacenterResource:
		return core.CheckRequiredGlobalFlags(c.Resource, config.ArgDataCenterId)
	case config.ServerResource:
		return core.CheckRequiredGlobalFlags(c.Resource, config.ArgDataCenterId, config.ArgServerId)
	case config.VolumeResource:
		return core.CheckRequiredGlobalFlags(c.Resource, config.ArgDataCenterId, config.ArgVolumeId)
	case config.IpBlockResource:
		return core.CheckRequiredGlobalFlags(c.Resource, config.ArgIpBlockId)
	case config.SnapshotResource:
		return core.CheckRequiredGlobalFlags(c.Resource, config.ArgSnapshotId)
	default:
		return c.Printer.Print(labelResourceWarning)
	}
}

// Output Printing

var defaultLabelCols = []string{"Key", "Value", "ResourceType", "ResourceId"}

type LabelPrint struct {
	Key          string `json:"Key,omitempty"`
	Value        string `json:"Value,omitempty"`
	ResourceType string `json:"ResourceType,omitempty"`
	ResourceId   string `json:"ResourceId,omitempty"`
}

func getLabelPrint(c *core.CommandConfig, s []v5.Label) printer.Result {
	r := printer.Result{}
	if c != nil {
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getLabelKVMaps(s)
			r.Columns = defaultLabelCols
		}
	}
	return r
}

func getLabels(Labels v5.Labels) []v5.Label {
	ss := make([]v5.Label, 0)
	if items, ok := Labels.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ss = append(ss, v5.Label{Label: s})
		}
	}
	return ss
}

func getLabel(s *v5.Label) []v5.Label {
	ss := make([]v5.Label, 0)
	if s != nil {
		ss = append(ss, v5.Label{Label: s.Label})
	}
	return ss
}

func getLabelKVMaps(ss []v5.Label) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		o := getLabelKVMap(s)
		out = append(out, o)
	}
	return out
}

func getLabelKVMap(s v5.Label) map[string]interface{} {
	var ssPrint LabelPrint
	if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
		if key, ok := properties.GetKeyOk(); ok && key != nil {
			ssPrint.Key = *key
		}
		if value, ok := properties.GetValueOk(); ok && value != nil {
			ssPrint.Value = *value
		}
		if resourceType, ok := properties.GetResourceTypeOk(); ok && resourceType != nil {
			ssPrint.ResourceType = *resourceType
		}
		if resourceId, ok := properties.GetResourceIdOk(); ok && resourceId != nil {
			ssPrint.ResourceId = *resourceId
		}
	}
	return structs.Map(ssPrint)
}
