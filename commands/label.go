package commands

import (
	"context"
	"os"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-multierror"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func label() *builder.Command {
	ctx := context.TODO()
	labelCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "label",
			Short:            "Label Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl label` + "`" + ` allow you to get, list, add, remove Labels from a Resource.`,
			TraverseChildren: true,
		},
	}
	globalFlags := labelCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultLabelResourceCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	globalFlags.StringP(config.ArgDataCenterId, "", "", "The unique Data Center Id")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelCmd.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgServerId, "", "", "The unique Server Id")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelCmd.Name(), config.ArgServerId), globalFlags.Lookup(config.ArgServerId))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(labelCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgVolumeId, "", "", "The unique Volume Id")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelCmd.Name(), config.ArgVolumeId), globalFlags.Lookup(config.ArgVolumeId))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(labelCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgIpBlockId, "", "", "The unique IpBlock Id")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelCmd.Name(), config.ArgIpBlockId), globalFlags.Lookup(config.ArgIpBlockId))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgSnapshotId, "", "", "The unique Snapshot Id")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelCmd.Name(), config.ArgSnapshotId), globalFlags.Lookup(config.ArgSnapshotId))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgResourceType, "", "", "Resource Type")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelCmd.Name(), config.ArgResourceType), globalFlags.Lookup(config.ArgResourceType))
	_ = labelCmd.Command.RegisterFlagCompletionFunc(config.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"datacenter", "volume", "server", "snapshot", "ipblock"}, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	builder.NewCommand(ctx, labelCmd, noPreRun, RunLabelList, "list", "List Labels from Resources",
		"Use this command to list all Labels from all Resources under your account. If you want to list all Labels from a specific Resource, use `--resource-type` option together with the Resource Id: `--datacenter-id`, `--server-id`, `--volume-id`.", listLabelsExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, labelCmd, PreRunGlobalResourceTypeLabelKey, RunLabelGet, "get", "Get a Label",
		"Use this command to get information about a specified Label from a specified Resource.\n\nRequired values to run command:\n\n* Resource Type\n* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id\n* Label Key",
		getLabelExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)

	/*
		Get By Urn Command
	*/
	getByUrn := builder.NewCommand(ctx, labelCmd, PreRunLabelUrn, RunLabelGetByUrn, "get-by-urn", "Get a Label using URN",
		"Use this command to get information about a specified Label using its URN. A URN is used for uniqueness of a Label and composed using `urn:label:<resource_type>:<resource_uuid>:<key>`.\n\nRequired values to run command:\n\n* Label URN",
		getLabelByUrnExample, true)
	getByUrn.AddStringFlag(config.ArgLabelUrn, "", "", "URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>] (required)")

	/*
		Add Command
	*/
	addLabel := builder.NewCommand(ctx, labelCmd, PreRunGlobalResourceTypeLabelKeyValue, RunLabelAdd, "add", "Add a Label to a Resource",
		"Use this command to add a Label to a specific Resource.\n\nRequired values to run command:\n\n* Resource Type\n* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id\n* Label Key\n* Label Value",
		addLabelExample, true)
	addLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	addLabel.AddStringFlag(config.ArgLabelValue, "", "", config.RequiredFlagLabelValue)

	/*
		Remove Command
	*/
	removeLabel := builder.NewCommand(ctx, labelCmd, PreRunGlobalResourceTypeLabelKey, RunLabelRemove, "remove", "Remove a Label from a Resource",
		"Use this command to remove a Label from a Resource.\n\nRequired values to run command:\n\n* Resource Type\n* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id\n* Label Key",
		removeLabelExample, true)
	removeLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)

	return labelCmd
}

func PreRunGlobalResourceTypeLabelKey(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgResourceType); err != nil {
		result = multierror.Append(result, err)
	}
	if err := checkResourceIds(c); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLabelKey); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func PreRunGlobalResourceTypeLabelKeyValue(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgResourceType); err != nil {
		result = multierror.Append(result, err)
	}
	if err := checkResourceIds(c); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLabelKey, config.ArgLabelValue); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func PreRunLabelUrn(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLabelUrn)
}

const labelResourceWarning = "Please use `--resource-type` flag with one option: \"datacenter\", \"volume\", \"server\", \"snapshot\", \"ipblock\""

func RunLabelList(c *builder.CommandConfig) error {
	switch viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgResourceType)) {
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

func RunLabelGet(c *builder.CommandConfig) error {
	switch viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgResourceType)) {
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

func RunLabelGetByUrn(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().GetByUrn(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelUrn)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelPrint(c, getLabel(labelDc)))
}

func RunLabelAdd(c *builder.CommandConfig) error {
	switch viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgResourceType)) {
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

func RunLabelRemove(c *builder.CommandConfig) error {
	switch viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgResourceType)) {
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

func checkResourceIds(c *builder.PreCommandConfig) error {
	switch viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgResourceType)) {
	case config.DatacenterResource:
		return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId)
	case config.ServerResource:
		return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId, config.ArgServerId)
	case config.VolumeResource:
		return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId, config.ArgVolumeId)
	case config.IpBlockResource:
		return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgIpBlockId)
	case config.SnapshotResource:
		return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgSnapshotId)
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

func getLabelPrint(c *builder.CommandConfig, s []resources.Label) printer.Result {
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
