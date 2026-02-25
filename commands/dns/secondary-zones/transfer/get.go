package transfer

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
			Aliases:   []string{"g"},
			ShortDesc: "Get the transfer status for a secondary zone",
			LongDesc:  "Get the transfer status for a secondary zone",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone)
			},
			CmdRun: func(c *core.CommandConfig) error {
				zoneNameOrID, _ := c.Command.Command.Flags().GetString(constants.FlagZone)
				zoneID, err := utils.SecondaryZoneResolve(zoneNameOrID)
				if err != nil {
					return err
				}

				transferStatuses, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesAxfrGet(context.Background(), zoneID).Execute()
				if err != nil {
					return err
				}

				cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				return table.Fprint(c.Command.Command.OutOrStdout(), allCols, transferStatuses, cols, table.WithPrefix("items"))
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
