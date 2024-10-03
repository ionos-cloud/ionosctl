package secondary_zones

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
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
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.SecondaryZonesIDs(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return c
}
