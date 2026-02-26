package secondary_zones

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

func getCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "get",
			ShortDesc: "Retrieve a secondary zone",
			LongDesc:  "Retrieve a secondary zone by its ID or name",
			Example:   "ionosctl dns secondary-zones get --zone ZONE_ID",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone)
			},
			CmdRun: func(c *core.CommandConfig) error {
				zoneNameOrID, _ := c.Command.Command.Flags().GetString(constants.FlagZone)
				zoneID, err := utils.SecondaryZoneResolve(zoneNameOrID)
				if err != nil {
					return err
				}

				secZone, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesFindById(context.Background(), zoneID).Execute()
				if err != nil {
					return err
				}

				cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				return c.Out(table.Sprint(allCols, secZone, cols))
			},
		},
	)

	c.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone,
		core.WithCompletion(completer.SecondaryZonesIDs, constants.DNSApiRegionalURL, constants.DNSLocations),
	)

	c.Command.SilenceUsage = true
	c.Command.Flags().SortFlags = false

	return c
}
