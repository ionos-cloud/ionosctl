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

func getRecordsPrint(c *core.CommandConfig, data ionoscloud.RecordsResponse) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil {
		// TODO for r.OutputJSON: This loses all kinds of information in `-o json`, like `limit`, `offset`, etc. See https://github.com/ionos-cloud/ionosctl/issues/249
		// But we are forced to do this otherwise we'd have this JSON output:
		// {
		//  "items": {
		//    "items": [
		// ...
		r.OutputJSON = data.Items // TODO: See above comment. Remove `.Items` once JSON marshalling works as one would expect
		r.KeyValue = makeRecordPrintObj(*data.Items...)
		r.Columns = printer.GetHeadersAllDefault(allCols, cols)
	}
	return r
}

func getRecordPrint(c *core.CommandConfig, data ionoscloud.RecordResponse) printer.Result {
	return getRecordsPrint(c, ionoscloud.RecordsResponse{Items: &[]ionoscloud.RecordResponse{data}})
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
