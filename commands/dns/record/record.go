package main

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func RecordsGetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "list",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve all records",
		Example:   "ionosctl dns record list ",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			/* TODO: Delete/modify me for --all
			 * err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.Flag<Parent>Id}, []string{constants.ArgAll, constants.Flag<Parent>Id})
			 * if err != nil {
			 * 	return err
			 * }
             * */

			// TODO: If no --all, mark individual flags as required

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			// Implement the actual command logic here
		},
		InitClient: true,
	})


	cmd.AddStringFlag(filter.zoneId, "", "", "Filter used to fetch only the records that contain specified zoneId")
	cmd.AddStringFlag(filter.name, "", "", "Filter used to fetch only the records that contain specified record name")
	cmd.AddIntFlag(offset, "", 0, "The first element (of the total list of elements) to include in the response. Use together with limit for pagination")
	cmd.AddIntFlag(limit, "", 0, "The maximum number of elements to return. Use together with offset for pagination")
	cmd.AddStringSliceFlag(constants.FlagItems, "", []string{}, "")
	cmd.AddFloat64Flag(constants.FlagLimit, "", 0.0, "Pagination limit")
	cmd.AddFloat64Flag(constants.FlagOffset, "", 0.0, "Pagination offset")

	return cmd
}

// Helper functions for printing record

func getrecordPrint(c *core.CommandConfig, dcs *[]ionoscloud.recordResponseData) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = makerecordPrintObj(dcs)
		r.Columns = printer.GetHeaders(allCols, defCols, cols)
	}
	return r
}

type recordPrint struct {
	Offset float `json:"Offset,omitempty"`
	Items []string `json:"Items,omitempty"`
	Limit float `json:"Limit,omitempty"`

}

var allCols = structs.Names(recordPrint{})
var defCols = allCols[:len(allCols)-3]

func makerecordPrintObj(data *[]ionoscloud.recordResponseData) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*data))

	for _, item := range *data {
		var printObj recordPrint
		printObj.Id = *item.GetId()

		// Fill in the rest of the fields from the response object

		if propertiesOk, ok := item.GetPropertiesOk(); ok && propertiesOk != nil {
			printObj.Offset = *propertiesOk.GetOffset()
		}
		if propertiesOk, ok := item.GetPropertiesOk(); ok && propertiesOk != nil {
			printObj.Items = *propertiesOk.GetItems()
		}
		if propertiesOk, ok := item.GetPropertiesOk(); ok && propertiesOk != nil {
			printObj.Limit = *propertiesOk.GetLimit()
		}

		o := structs.Map(printObj)
		out = append(out, o)
	}
	return out
}
