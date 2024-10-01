package secondary_zones

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
)

func listCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "list",
			ShortDesc: "List secondary zones",
			LongDesc:  "List secondary zones",
			Example:   "ionosctl dns secondary-zone list",
			PreCmdRun: nil,
			CmdRun: func(c *core.CommandConfig) error {
				secZones, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesGet(context.Background()).Execute()
				if err != nil {
					return err
				}

				cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				out, err := jsontabwriter.GenerateOutput(
					"items", jsonpaths.DnsSecondaryZone, secZones, tabheaders.GetHeadersAllDefault(allCols, cols),
				)
				if err != nil {
					return err
				}

				fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
				return nil
			},
		},
	)

	return c
}
