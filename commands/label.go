package commands

import (
	"context"

	"github.com/fatih/structs"
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
			Long:             `The sub-commands of ` + "`" + `ionosctl label` + "`" + ` allow you to get, list, create, delete Labels from a resource. For each resource that supports labelling: Data Center, Server, Volume, IpBlock, Snapshot - commands to manage Labels are available. Example: ` + "`" + `ionosctl <resource> add-label` + "`" + `.`,
			TraverseChildren: true,
		},
	}
	globalFlags := labelCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultLabelCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, labelCmd, noPreRun, RunLabelList, "list", "List Labels from all Resources",
		"Use this command to list all Labels from all Resources under your account.", listLabelsExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, labelCmd, PreRunLabelUrn, RunLabelGet, "get", "Get a Label",
		"Use this command to get information about a specified Label using its URN. A URN is used for uniqueness of a Label and composed using `urn:label:<resource_type>:<resource_uuid>:<key>`.\n\nRequired values to run command:\n\n* Label URN",
		getLabelExample, true)
	get.AddStringFlag(config.ArgLabelUrn, "", "", "URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>]")

	return labelCmd
}

func PreRunLabelUrn(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLabelUrn)
}

func RunLabelList(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelPrint(c, getLabels(labelDcs)))
}

func RunLabelGet(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().GetByUrn(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelUrn)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelPrint(c, getLabel(labelDc)))
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
