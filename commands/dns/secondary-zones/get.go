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

func getCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "get",
			ShortDesc: "Retrieve a secondary zone",
			LongDesc:  "Retrieve a secondary zone",
			Example:   "ionosctl dns secondary-zones get --zone ZONE_ID",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone)
			},
			CmdRun: func(c *core.CommandConfig) error {
				zone, _ := c.Command.Command.Flags().GetString(constants.FlagZone)
				zoneID, err := resolve(zone)
				if err != nil {
					return err
				}

				secZone, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesFindById(context.Background(), zoneID).Execute()
				if err != nil {
					return err
				}

				cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				out, err := jsontabwriter.GenerateOutput(
					"", jsonpaths.DnsSecondaryZone, secZone, tabheaders.GetHeadersAllDefault(allCols, cols),
				)
				if err != nil {
					return err
				}

				fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
				return nil
			},
		},
	)

	c.Command.Flags().StringP(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone)

	return c
}
