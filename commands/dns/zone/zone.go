package zone

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dnsaas"
	"github.com/spf13/cobra"
)

func ZoneCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "zone",
			Short:            "DNS zones",
			Aliases:          []string{"z"},
			Long:             "The sub-commands of `ionosctl dns zone` allow you to perform operations on DNS zones",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(ZonesGetCmd())
	cmd.AddCommand(ZonesDeleteCmd())
	cmd.AddCommand(ZonesPostCmd())
	cmd.AddCommand(ZonesPutCmd())
	cmd.AddCommand(ZonesFindByIdCmd())
	return cmd
}

// Helper functions for printing zone

func getZonesPrint(c *core.CommandConfig, data ionoscloud.ZonesResponse) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil {
		r.OutputJSON = data
		r.KeyValue = makeZonePrintObj(*data.Items...)
		r.Columns = printer.GetHeadersAllDefault(allCols, cols)
	}
	return r
}

func getZonePrint(c *core.CommandConfig, data ionoscloud.ZoneResponse) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil {
		r.OutputJSON = data
		r.KeyValue = makeZonePrintObj(data)
		r.Columns = printer.GetHeadersAllDefault(allCols, cols)
	}
	return r
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
