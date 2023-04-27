package record

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/spf13/cobra"
)

func RecordCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "record",
			Short:            "DNS Records",
			Aliases:          []string{"r"},
			Long:             "The sub-commands of `ionosctl record` allow you to perform operations on DBaaS resources.",
			TraverseChildren: true,
		},
	}
	dbaasCmd.AddCommand(postgres.DBaaSPostgresCmd())
	dbaasCmd.AddCommand(mongo.DBaaSMongoCmd())
	return dbaasCmd
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
	Offset float    `json:"Offset,omitempty"`
	Items  []string `json:"Items,omitempty"`
	Limit  float    `json:"Limit,omitempty"`
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
