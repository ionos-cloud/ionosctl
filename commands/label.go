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

	/*
		List Command
	*/
	list := core.NewCommand(ctx, labelCmd, core.CommandBuilder{
		Namespace:  "label",
		Resource:   "label",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Labels from Resources",
		LongDesc:   "Use this command to list all Labels from all Resources under your account. If you want to list all Labels from a specific Resource, use `--resource-type` option together with the Resource Id: `--datacenter-id`, `--server-id`, `--volume-id`.",
		Example:    listLabelsExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunLabelList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgServerId, "", "", config.ServerId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgVolumeId, "", "", config.VolumeId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgIpBlockId, "", "", config.IpBlockId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgSnapshotId, "", "", config.SnapshotId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgResourceType, "", "", "Resource Type")
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{config.DatacenterResource, config.VolumeResource, config.ServerResource, config.SnapshotResource, config.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(config.ArgLabelKey, "", "", config.LabelKey, core.RequiredFlagOption())
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgServerId, "", "", config.ServerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgVolumeId, "", "", config.VolumeId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgIpBlockId, "", "", config.IpBlockId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgSnapshotId, "", "", config.SnapshotId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgResourceType, "", "", "Resource Type", core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{config.DatacenterResource, config.VolumeResource, config.ServerResource, config.SnapshotResource, config.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
	})

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
	getByUrn.AddStringFlag(config.ArgLabelUrn, "", "", "URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>]", core.RequiredFlagOption())

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
	addLabel.AddStringFlag(config.ArgLabelKey, "", "", config.LabelKey, core.RequiredFlagOption())
	addLabel.AddStringFlag(config.ArgLabelValue, "", "", config.LabelValue, core.RequiredFlagOption())
	addLabel.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(config.ArgServerId, "", "", config.ServerId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(addLabel.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(config.ArgVolumeId, "", "", config.VolumeId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(addLabel.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(config.ArgIpBlockId, "", "", config.IpBlockId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(config.ArgSnapshotId, "", "", config.SnapshotId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(config.ArgResourceType, "", "", "Resource Type", core.RequiredFlagOption())
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{config.DatacenterResource, config.VolumeResource, config.ServerResource, config.SnapshotResource, config.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
	})

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
		PreCmdRun:  PreRunResourceTypeLabelKey,
		CmdRun:     RunLabelRemove,
		InitClient: true,
	})
	removeLabel.AddStringFlag(config.ArgLabelKey, "", "", config.LabelKey, core.RequiredFlagOption())
	removeLabel.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(config.ArgServerId, "", "", config.ServerId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(removeLabel.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(config.ArgVolumeId, "", "", config.VolumeId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(removeLabel.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(config.ArgIpBlockId, "", "", config.IpBlockId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(config.ArgSnapshotId, "", "", config.SnapshotId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(config.ArgResourceType, "", "", "Resource Type", core.RequiredFlagOption())
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{config.DatacenterResource, config.VolumeResource, config.ServerResource, config.SnapshotResource, config.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
	})

	return labelCmd
}

func PreRunResourceTypeLabelKey(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{config.ArgResourceType, config.ArgDataCenterId, config.ArgLabelKey},
		[]string{config.ArgResourceType, config.ArgDataCenterId, config.ArgVolumeId, config.ArgLabelKey},
		[]string{config.ArgResourceType, config.ArgDataCenterId, config.ArgServerId, config.ArgLabelKey},
		[]string{config.ArgResourceType, config.ArgSnapshotId, config.ArgLabelKey},
		[]string{config.ArgResourceType, config.ArgIpBlockId, config.ArgLabelKey},
	)
}

func PreRunResourceTypeLabelKeyValue(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{config.ArgResourceType, config.ArgDataCenterId, config.ArgLabelKey, config.ArgLabelValue},
		[]string{config.ArgResourceType, config.ArgDataCenterId, config.ArgVolumeId, config.ArgLabelKey, config.ArgLabelValue},
		[]string{config.ArgResourceType, config.ArgDataCenterId, config.ArgServerId, config.ArgLabelKey, config.ArgLabelValue},
		[]string{config.ArgResourceType, config.ArgSnapshotId, config.ArgLabelKey, config.ArgLabelValue},
		[]string{config.ArgResourceType, config.ArgIpBlockId, config.ArgLabelKey, config.ArgLabelValue},
	)
}

func PreRunLabelUrn(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, config.ArgLabelUrn)
}

const labelResourceWarning = "Please use `--resource-type` flag with one option: \"datacenter\", \"volume\", \"server\", \"snapshot\", \"ipblock\""

func RunLabelList(c *core.CommandConfig) error {
	switch viper.GetString(core.GetFlagName(c.NS, config.ArgResourceType)) {
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
	resourceType := viper.GetString(core.GetFlagName(c.NS, config.ArgResourceType))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Getting label with label key: %v and label value: %v for %v...", labelKey, labelValue, resourceType)
	switch resourceType {
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
	urn := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelUrn))
	c.Printer.Verbose("Getting label with urn: %v", urn)
	labelDc, _, err := c.Labels().GetByUrn(urn)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelPrint(c, getLabel(labelDc)))
}

func RunLabelAdd(c *core.CommandConfig) error {
	switch viper.GetString(core.GetFlagName(c.NS, config.ArgResourceType)) {
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
	resourceType := viper.GetString(core.GetFlagName(c.NS, config.ArgResourceType))
	c.Printer.Verbose("Label is removing from %v...", resourceType)
	switch resourceType {
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
