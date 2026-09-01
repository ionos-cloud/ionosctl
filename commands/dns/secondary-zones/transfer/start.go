package transfer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func startCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "start",
			Aliases:   []string{"s"},
			ShortDesc: "Start a zone transfer from the primary",
			LongDesc:  "Trigger IONOS to pull the latest records for a secondary zone from its primary name servers. Run this after creating the zone or whenever the primary changes; check progress with 'transfer get'.",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone)
			},
			CmdRun: func(c *core.CommandConfig) error {
				zoneNameOrID, _ := c.Command.Command.Flags().GetString(constants.FlagZone)
				zoneID, err := utils.SecondaryZoneResolve(zoneNameOrID)
				if err != nil {
					return err
				}

				_, _, err = client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesAxfrPut(context.Background(), zoneID).Execute()
				if err != nil {
					return err
				}

				c.Msg("Transfer started for zone %s", zoneID)
				return nil
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
