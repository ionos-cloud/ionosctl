package transfer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/spf13/cobra"
)

func startCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "start",
			Aliases:   []string{"s"},
			ShortDesc: "Initiate zone transfer",
			LongDesc:  "Initiate zone transfer",
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

				fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Transfer started for zone %s", zoneID))
				return nil
			},
		},
	)

	c.Command.Flags().StringP(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone)
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.SecondaryZonesIDs(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.Command.SilenceUsage = true
	c.Command.Flags().SortFlags = false

	return c
}
