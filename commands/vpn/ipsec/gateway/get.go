package gateway

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "ipsec gateway",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Find a gateway by ID",
		Example:   "ionosctl vpn ipsec gateway get " + core.FlagsUsage(constants.FlagGatewayID),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID)
		},
		CmdRun: func(c *core.CommandConfig) error {
			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))

			g, _, err := client.Must().VPNClient.IPSecGatewaysApi.IpsecgatewaysFindById(context.Background(), id).Execute()
			if err != nil {
				return fmt.Errorf("failed getting gateway by id %s: %w", id, err)
			}

			table, err := resource2table.ConvertVPNIPSecGatewayToTable(g)
			if err != nil {
				return fmt.Errorf("could not convert from JSON to Table format: %w", err)
			}
			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutputPreconverted(g, table,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the IPSec Gateway", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagGatewayID, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GatewayIDs(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
