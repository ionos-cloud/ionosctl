package gateway

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "ipsec gateway",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List IPSec Gateways",
		Example:   "ionosctl vpn ipsec gateway list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, err := completer.Gateways(
				func(req vpn.ApiIpsecgatewaysGetRequest) (vpn.ApiIpsecgatewaysGetRequest, error) {
					if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
						req = req.Offset(viper.GetInt32(fn))
					}
					if fn := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(fn) {
						req = req.Limit(viper.GetInt32(fn))
					}
					return req, nil
				},
			)
			if err != nil {
				return fmt.Errorf("failed listing gateways: %w", err)
			}

			table, err := resource2table.ConvertVPNIPSecGatewaysToTable(ls)
			if err != nil {
				return fmt.Errorf("could not convert from JSON to Table format: %w", err)
			}
			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutputPreconverted(ls, table,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
	})

	cmd.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, 0, constants.DescMaxResults)
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "Skip a certain number of results")

	return cmd
}
