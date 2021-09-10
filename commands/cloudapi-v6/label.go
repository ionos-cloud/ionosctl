package commands

import (
	"context"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
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
	list.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapi_v6.ArgVolumeId, "", "", cloudapi_v6.VolumeId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapi_v6.ArgIpBlockId, "", "", cloudapi_v6.IpBlockId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapi_v6.ArgSnapshotId, "", "", cloudapi_v6.SnapshotId)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapi_v6.ArgResourceType, "", "", "Type of the resource to list labels from")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{cloudapi_v6.DatacenterResource, cloudapi_v6.VolumeResource, cloudapi_v6.ServerResource, cloudapi_v6.SnapshotResource, cloudapi_v6.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(cloudapi_v6.ArgLabelKey, "", "", cloudapi_v6.LabelKey, core.RequiredFlagOption())
	get.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgVolumeId, "", "", cloudapi_v6.VolumeId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgIpBlockId, "", "", cloudapi_v6.IpBlockId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgSnapshotId, "", "", cloudapi_v6.SnapshotId)
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgResourceType, "", "", "Type of the resource to get label from", core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{cloudapi_v6.DatacenterResource, cloudapi_v6.VolumeResource, cloudapi_v6.ServerResource, cloudapi_v6.SnapshotResource, cloudapi_v6.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
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
	getByUrn.AddStringFlag(cloudapi_v6.ArgLabelUrn, "", "", "URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>]", core.RequiredFlagOption())

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
	addLabel.AddStringFlag(cloudapi_v6.ArgLabelKey, "", "", cloudapi_v6.LabelKey, core.RequiredFlagOption())
	addLabel.AddStringFlag(cloudapi_v6.ArgLabelValue, "", "", cloudapi_v6.LabelValue, core.RequiredFlagOption())
	addLabel.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(addLabel.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(cloudapi_v6.ArgVolumeId, "", "", cloudapi_v6.VolumeId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(addLabel.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(cloudapi_v6.ArgIpBlockId, "", "", cloudapi_v6.IpBlockId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(cloudapi_v6.ArgSnapshotId, "", "", cloudapi_v6.SnapshotId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addLabel.AddStringFlag(cloudapi_v6.ArgResourceType, "", "", "Type of the resource to add label to", core.RequiredFlagOption())
	_ = addLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{cloudapi_v6.DatacenterResource, cloudapi_v6.VolumeResource, cloudapi_v6.ServerResource, cloudapi_v6.SnapshotResource, cloudapi_v6.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
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
	removeLabel.AddStringFlag(cloudapi_v6.ArgLabelKey, "", "", cloudapi_v6.LabelKey, core.RequiredFlagOption())
	removeLabel.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(removeLabel.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(cloudapi_v6.ArgVolumeId, "", "", cloudapi_v6.VolumeId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(removeLabel.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(cloudapi_v6.ArgIpBlockId, "", "", cloudapi_v6.IpBlockId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(cloudapi_v6.ArgSnapshotId, "", "", cloudapi_v6.SnapshotId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(cloudapi_v6.ArgResourceType, "", "", "Type of the resource to remove label for", core.RequiredFlagOption())
	_ = removeLabel.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{cloudapi_v6.DatacenterResource, cloudapi_v6.VolumeResource, cloudapi_v6.ServerResource, cloudapi_v6.SnapshotResource, cloudapi_v6.IpBlockResource}, cobra.ShellCompDirectiveNoFileComp
	})

	return labelCmd
}

func PreRunResourceTypeLabelKey(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapi_v6.ArgResourceType, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgLabelKey},
		[]string{cloudapi_v6.ArgResourceType, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgVolumeId, cloudapi_v6.ArgLabelKey},
		[]string{cloudapi_v6.ArgResourceType, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgServerId, cloudapi_v6.ArgLabelKey},
		[]string{cloudapi_v6.ArgResourceType, cloudapi_v6.ArgSnapshotId, cloudapi_v6.ArgLabelKey},
		[]string{cloudapi_v6.ArgResourceType, cloudapi_v6.ArgIpBlockId, cloudapi_v6.ArgLabelKey},
	)
}

func PreRunResourceTypeLabelKeyValue(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapi_v6.ArgResourceType, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgLabelKey, cloudapi_v6.ArgLabelValue},
		[]string{cloudapi_v6.ArgResourceType, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgVolumeId, cloudapi_v6.ArgLabelKey, cloudapi_v6.ArgLabelValue},
		[]string{cloudapi_v6.ArgResourceType, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgServerId, cloudapi_v6.ArgLabelKey, cloudapi_v6.ArgLabelValue},
		[]string{cloudapi_v6.ArgResourceType, cloudapi_v6.ArgSnapshotId, cloudapi_v6.ArgLabelKey, cloudapi_v6.ArgLabelValue},
		[]string{cloudapi_v6.ArgResourceType, cloudapi_v6.ArgIpBlockId, cloudapi_v6.ArgLabelKey, cloudapi_v6.ArgLabelValue},
	)
}

func PreRunLabelUrn(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapi_v6.ArgLabelUrn)
}

const labelResourceWarning = "Please use `--resource-type` flag with one option: \"datacenter\", \"volume\", \"server\", \"snapshot\", \"ipblock\""

func RunLabelList(c *core.CommandConfig) error {
	switch viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgResourceType)) {
	case cloudapi_v6.DatacenterResource:
		return RunDataCenterLabelsList(c)
	case cloudapi_v6.ServerResource:
		return RunServerLabelsList(c)
	case cloudapi_v6.VolumeResource:
		return RunVolumeLabelsList(c)
	case cloudapi_v6.IpBlockResource:
		return RunIpBlockLabelsList(c)
	case cloudapi_v6.SnapshotResource:
		return RunSnapshotLabelsList(c)
	default:
		labelDcs, _, err := c.CloudApiV6Services.Labels().List()
		if err != nil {
			return err
		}
		return c.Printer.Print(getLabelPrint(c, getLabels(labelDcs)))
	}
}

func RunLabelGet(c *core.CommandConfig) error {
	resourceType := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgResourceType))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgLabelValue))
	c.Printer.Verbose("Getting label with label key: %v and label value: %v for %v...", labelKey, labelValue, resourceType)
	switch resourceType {
	case cloudapi_v6.DatacenterResource:
		return RunDataCenterLabelGet(c)
	case cloudapi_v6.ServerResource:
		return RunServerLabelGet(c)
	case cloudapi_v6.VolumeResource:
		return RunVolumeLabelGet(c)
	case cloudapi_v6.IpBlockResource:
		return RunIpBlockLabelGet(c)
	case cloudapi_v6.SnapshotResource:
		return RunSnapshotLabelGet(c)
	default:
		return c.Printer.Print(labelResourceWarning)
	}
}

func RunLabelGetByUrn(c *core.CommandConfig) error {
	urn := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgLabelUrn))
	c.Printer.Verbose("Getting label with urn: %v", urn)
	labelDc, _, err := c.CloudApiV6Services.Labels().GetByUrn(urn)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelPrint(c, getLabel(labelDc)))
}

func RunLabelAdd(c *core.CommandConfig) error {
	switch viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgResourceType)) {
	case cloudapi_v6.DatacenterResource:
		return RunDataCenterLabelAdd(c)
	case cloudapi_v6.ServerResource:
		return RunServerLabelAdd(c)
	case cloudapi_v6.VolumeResource:
		return RunVolumeLabelAdd(c)
	case cloudapi_v6.IpBlockResource:
		return RunIpBlockLabelAdd(c)
	case cloudapi_v6.SnapshotResource:
		return RunSnapshotLabelAdd(c)
	default:
		return c.Printer.Print(labelResourceWarning)
	}
}

func RunLabelRemove(c *core.CommandConfig) error {
	resourceType := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgResourceType))
	c.Printer.Verbose("Label is removing from %v...", resourceType)
	switch resourceType {
	case cloudapi_v6.DatacenterResource:
		return RunDataCenterLabelRemove(c)
	case cloudapi_v6.ServerResource:
		return RunServerLabelRemove(c)
	case cloudapi_v6.VolumeResource:
		return RunVolumeLabelRemove(c)
	case cloudapi_v6.IpBlockResource:
		return RunIpBlockLabelRemove(c)
	case cloudapi_v6.SnapshotResource:
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
