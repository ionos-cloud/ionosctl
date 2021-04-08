package commands

import (
	"context"
	"errors"
	"io"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
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
			Long:             `The sub-commands of ` + "`" + `ionosctl label` + "`" + ` allow you to get, list, create, delete Labels from a Resource`,
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	builder.NewCommand(ctx, labelCmd, noPreRun, RunLabelList, "list", "List Labels from a Data Center",
		"Use this command to list all Labels from a specified Data Center.", listDatacenterExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, labelCmd, noPreRun, RunLabelGet, "get", "Get a Label from a Data Center",
		"Use this command to get information about a specified Label from a Data Center.\n\nRequired values to run command:\n\n* Data Center Id\n* Label Key", getDatacenterExample, true)
	get.AddStringFlag(config.ArgLabelUrn, "", "", "")

	labelCmd.AddCommand(labelDatacenter())
	labelCmd.AddCommand(labelServer())
	labelCmd.AddCommand(labelVolume())
	labelCmd.AddCommand(labelIpBlock())
	labelCmd.AddCommand(labelSnapshot())
	return labelCmd
}

func RunLabelList(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelPrint(nil, c, getLabels(labelDcs)))
}

func RunLabelGet(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().GetByUrn(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelUrn)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelPrint(nil, c, getLabel(labelDc)))
}

// Output Printing

var defaultLabelCols = []string{"Key", "Value", "ResourceType", "ResourceId"}

type LabelPrint struct {
	Key         string `json:"Key,omitempty"`
	Value       string `json:"Value,omitempty"`
	ResourceType string `json:"ResourceType,omitempty"`
	ResourceId  string `json:"ResourceId,omitempty"`
}

func getLabelPrint(resp *resources.Response, c *builder.CommandConfig, s []resources.Label) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = "label " + c.ParentName
			r.Verb = c.Name
		}
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getLabelKVMaps(s)
			r.Columns = getLabelCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getLabelCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var labelCols []string
		columnsMap := map[string]string{
			"Key":          "Key",
			"Value":        "Value",
			"ResourceType": "ResourceType",
			"ResourceId":   "ResourceId",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				labelCols = append(labelCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return labelCols
	} else {
		return defaultLabelCols
	}
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
