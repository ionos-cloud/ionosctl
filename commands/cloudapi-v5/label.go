package cloudapi_v5

import (
	"context"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/pkg/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/pkg/cloudapi-v5/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultLabelResourceCols, printer.ColsMessage(defaultLabelResourceCols))
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
	list.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv5.ArgServerId, "", "", cloudapiv5.ServerId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv5.ArgVolumeId, "", "", cloudapiv5.VolumeId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv5.ArgIpBlockId, "", "", cloudapiv5.IpBlockId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv5.ArgSnapshotId, "", "", cloudapiv5.SnapshotId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv5.ArgResourceType, "", "", "Type of the resource to list labels from")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{cloudapiv5.DatacenterResource, cloudapiv5.VolumeResource, cloudapiv5.ServerResource, cloudapiv5.SnapshotResource, cloudapiv5.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(cloudapiv5.ArgLabelKey, "", "", cloudapiv5.LabelKey, core.RequiredFlagOption())
	get.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgServerId, "", "", cloudapiv5.ServerId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgVolumeId, "", "", cloudapiv5.VolumeId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgIpBlockId, "", "", cloudapiv5.IpBlockId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgSnapshotId, "", "", cloudapiv5.SnapshotId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgResourceType, "", "", "Type of the resource to get label from", core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{cloudapiv5.DatacenterResource, cloudapiv5.VolumeResource, cloudapiv5.ServerResource, cloudapiv5.SnapshotResource, cloudapiv5.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
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
	getByUrn.AddStringFlag(cloudapiv5.ArgLabelUrn, "", "", "URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>]", core.RequiredFlagOption())

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
	addLabel.AddStringFlag(cloudapiv5.ArgLabelKey, "", "", cloudapiv5.LabelKey, core.RequiredFlagOption())
	addLabel.AddStringFlag(cloudapiv5.ArgLabelValue, "", "", cloudapiv5.LabelValue, core.RequiredFlagOption())
	addLabel.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(cloudapiv5.ArgServerId, "", "", cloudapiv5.ServerId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(addLabel.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(cloudapiv5.ArgVolumeId, "", "", cloudapiv5.VolumeId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(addLabel.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(cloudapiv5.ArgIpBlockId, "", "", cloudapiv5.IpBlockId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(cloudapiv5.ArgSnapshotId, "", "", cloudapiv5.SnapshotId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(cloudapiv5.ArgResourceType, "", "", "Type of the resource to add label to", core.RequiredFlagOption())
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{cloudapiv5.DatacenterResource, cloudapiv5.VolumeResource, cloudapiv5.ServerResource, cloudapiv5.SnapshotResource, cloudapiv5.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
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
	removeLabel.AddStringFlag(cloudapiv5.ArgLabelKey, "", "", cloudapiv5.LabelKey, core.RequiredFlagOption())
	removeLabel.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(cloudapiv5.ArgServerId, "", "", cloudapiv5.ServerId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(removeLabel.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(cloudapiv5.ArgVolumeId, "", "", cloudapiv5.VolumeId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(removeLabel.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(cloudapiv5.ArgIpBlockId, "", "", cloudapiv5.IpBlockId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(cloudapiv5.ArgSnapshotId, "", "", cloudapiv5.SnapshotId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(cloudapiv5.ArgResourceType, "", "", "Type of the resource to remove label for", core.RequiredFlagOption())
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{cloudapiv5.DatacenterResource, cloudapiv5.VolumeResource, cloudapiv5.ServerResource, cloudapiv5.SnapshotResource, cloudapiv5.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
	})

	return labelCmd
}

func PreRunResourceTypeLabelKey(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv5.ArgResourceType, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgLabelKey},
		[]string{cloudapiv5.ArgResourceType, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgVolumeId, cloudapiv5.ArgLabelKey},
		[]string{cloudapiv5.ArgResourceType, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgServerId, cloudapiv5.ArgLabelKey},
		[]string{cloudapiv5.ArgResourceType, cloudapiv5.ArgSnapshotId, cloudapiv5.ArgLabelKey},
		[]string{cloudapiv5.ArgResourceType, cloudapiv5.ArgIpBlockId, cloudapiv5.ArgLabelKey},
	)
}

func PreRunResourceTypeLabelKeyValue(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv5.ArgResourceType, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgLabelKey, cloudapiv5.ArgLabelValue},
		[]string{cloudapiv5.ArgResourceType, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgVolumeId, cloudapiv5.ArgLabelKey, cloudapiv5.ArgLabelValue},
		[]string{cloudapiv5.ArgResourceType, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgServerId, cloudapiv5.ArgLabelKey, cloudapiv5.ArgLabelValue},
		[]string{cloudapiv5.ArgResourceType, cloudapiv5.ArgSnapshotId, cloudapiv5.ArgLabelKey, cloudapiv5.ArgLabelValue},
		[]string{cloudapiv5.ArgResourceType, cloudapiv5.ArgIpBlockId, cloudapiv5.ArgLabelKey, cloudapiv5.ArgLabelValue},
	)
}

func PreRunLabelUrn(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgLabelUrn)
}

const labelResourceWarning = "Please use `--resource-type` flag with one option: \"datacenter\", \"volume\", \"server\", \"snapshot\", \"ipblock\""

func RunLabelList(c *core.CommandConfig) error {
	switch viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgResourceType)) {
	case cloudapiv5.DatacenterResource:
		return RunDataCenterLabelsList(c)
	case cloudapiv5.ServerResource:
		return RunServerLabelsList(c)
	case cloudapiv5.VolumeResource:
		return RunVolumeLabelsList(c)
	case cloudapiv5.IpBlockResource:
		return RunIpBlockLabelsList(c)
	case cloudapiv5.SnapshotResource:
		return RunSnapshotLabelsList(c)
	default:
		labelDcs, _, err := c.CloudApiV5Services.Labels().List()
		if err != nil {
			return err
		}
		return c.Printer.Print(getLabelPrint(c, getLabels(labelDcs)))
	}
}

func RunLabelGet(c *core.CommandConfig) error {
	resourceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgResourceType))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelValue))
	c.Printer.Verbose("Getting label with label key: %v and label value: %v for %v...", labelKey, labelValue, resourceType)
	switch resourceType {
	case cloudapiv5.DatacenterResource:
		return RunDataCenterLabelGet(c)
	case cloudapiv5.ServerResource:
		return RunServerLabelGet(c)
	case cloudapiv5.VolumeResource:
		return RunVolumeLabelGet(c)
	case cloudapiv5.IpBlockResource:
		return RunIpBlockLabelGet(c)
	case cloudapiv5.SnapshotResource:
		return RunSnapshotLabelGet(c)
	default:
		return c.Printer.Print(labelResourceWarning)
	}
}

func RunLabelGetByUrn(c *core.CommandConfig) error {
	urn := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelUrn))
	c.Printer.Verbose("Getting label with urn: %v", urn)
	labelDc, _, err := c.CloudApiV5Services.Labels().GetByUrn(urn)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelPrint(c, getLabel(labelDc)))
}

func RunLabelAdd(c *core.CommandConfig) error {
	switch viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgResourceType)) {
	case cloudapiv5.DatacenterResource:
		return RunDataCenterLabelAdd(c)
	case cloudapiv5.ServerResource:
		return RunServerLabelAdd(c)
	case cloudapiv5.VolumeResource:
		return RunVolumeLabelAdd(c)
	case cloudapiv5.IpBlockResource:
		return RunIpBlockLabelAdd(c)
	case cloudapiv5.SnapshotResource:
		return RunSnapshotLabelAdd(c)
	default:
		return c.Printer.Print(labelResourceWarning)
	}
}

func RunLabelRemove(c *core.CommandConfig) error {
	resourceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgResourceType))
	c.Printer.Verbose("Label is removing from %v...", resourceType)
	switch resourceType {
	case cloudapiv5.DatacenterResource:
		return RunDataCenterLabelRemove(c)
	case cloudapiv5.ServerResource:
		return RunServerLabelRemove(c)
	case cloudapiv5.VolumeResource:
		return RunVolumeLabelRemove(c)
	case cloudapiv5.IpBlockResource:
		return RunIpBlockLabelRemove(c)
	case cloudapiv5.SnapshotResource:
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

func getLabelPrint(c *core.CommandConfig, s []resources.Label) printer.Result {
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

func getLabels(Labels resources.Labels) []resources.Label {
	ss := make([]resources.Label, 0)
	if items, ok := Labels.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ss = append(ss, resources.Label{Label: s})
		}
	}
	return ss
}

func getLabel(s *resources.Label) []resources.Label {
	ss := make([]resources.Label, 0)
	if s != nil {
		ss = append(ss, resources.Label{Label: s.Label})
	}
	return ss
}

func getLabelKVMaps(ss []resources.Label) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		o := getLabelKVMap(s)
		out = append(out, o)
	}
	return out
}

func getLabelKVMap(s resources.Label) map[string]interface{} {
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
