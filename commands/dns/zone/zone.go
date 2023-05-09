package zone

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/spf13/cobra"
)

func ZoneCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "zone",
			Short:            "DNS zones",
			Aliases:          []string{"z", "zones"},
			Long:             "The sub-commands of `ionosctl dns zone` allow you to perform operations on DNS zones",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(ZonesGetCmd())
	cmd.AddCommand(ZonesDeleteCmd())
	cmd.AddCommand(ZonesPostCmd())
	cmd.AddCommand(ZonesPutCmd())
	cmd.AddCommand(ZonesFindByIdCmd())

	// Quality-Of-Life commands which use another command in their implementation
	cmd.AddCommand(ZonesEnableCmd())
	cmd.AddCommand(ZonesDisableCmd())

	return cmd
}

// Helper functions for printing zone

func getZonesPrint(c *core.CommandConfig, data ionoscloud.ZonesResponse) printer.Result {
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
		r.KeyValue = makeZonePrintObj(*data.Items...)
		r.Columns = printer.GetHeadersAllDefault(allCols, cols)
	}
	return r
}

func getZonePrint(c *core.CommandConfig, data ionoscloud.ZoneResponse) printer.Result {
	return getZonesPrint(c, ionoscloud.ZonesResponse{Items: &[]ionoscloud.ZoneResponse{data}})
}

type zonePrint struct {
	Id          string `json:"ID,omitempty"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Content,omitempty"`
	Enabled     bool   `json:"Enabled,omitempty"`
	State       string `json:"State,omitempty"`
}

var allCols = structs.Names(zonePrint{})

func makeZonePrintObj(data ...ionoscloud.ZoneResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(data))

	for _, item := range data {
		var printObj zonePrint
		printObj.Id = *item.GetId()

		if p, ok := item.GetPropertiesOk(); ok {
			printObj.Enabled = *p.Enabled
			printObj.Description = *p.Description
			printObj.Name = *p.ZoneName
		}
		if m, ok := item.GetMetadataOk(); ok && m != nil {
			printObj.State = string(*m.State)
		}

		o := structs.Map(printObj)
		out = append(out, o)
	}
	return out
}
