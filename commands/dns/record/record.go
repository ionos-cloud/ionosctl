package record

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dnsaas"
	"github.com/spf13/cobra"
)

func RecordCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "record",
			Short:            "DNS Records",
			Aliases:          []string{"r"},
			Long:             "The sub-commands of `ionosctl dns record` allow you to perform operations on DNS records",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(RecordsGetCmd())
	cmd.AddCommand(ZonesRecordsDeleteCmd())
	cmd.AddCommand(ZonesRecordsPostCmd())
	cmd.AddCommand(ZonesRecordsFindByIdCmd())
	cmd.AddCommand(ZonesRecordsPutCmd())
	return cmd
}

// Helper functions for printing record

func getRecordPrint(c *core.CommandConfig, data ionoscloud.RecordsResponse) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil {
		r.OutputJSON = data
		r.KeyValue = makeRecordPrintObj(*data.Items...)
		r.Columns = printer.GetHeadersAllDefault(allCols, cols)
	}
	return r
}

type recordPrint struct {
	Id      string `json:"ID,omitempty"`
	Name    string `json:"Name,omitempty"`
	Content string `json:"Content,omitempty"`
	Type    string `json:"Type,omitempty"`
	Enabled bool   `json:"Enabled,omitempty"`
	State   string `json:"State,omitempty"`
}

var allCols = structs.Names(recordPrint{})

func makeRecordPrintObj(data ...ionoscloud.RecordResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(data))

	for _, item := range data {
		var printObj recordPrint
		printObj.Id = *item.GetId()

		// Fill in the rest of the fields from the response object

		if propertiesOk, ok := item.GetPropertiesOk(); ok && propertiesOk != nil {
			var j []byte
			err := propertiesOk.Type.UnmarshalJSON(j)
			if err == nil {
				printObj.Type = string(j)
			}

			printObj.Enabled = *propertiesOk.Enabled
			printObj.Content = *propertiesOk.Content
			printObj.Name = *propertiesOk.Name
		}
		if m, ok := item.GetMetadataOk(); ok && m != nil {
			printObj.State = string(*m.State)
		}

		o := structs.Map(printObj)
		out = append(out, o)
	}
	return out
}
