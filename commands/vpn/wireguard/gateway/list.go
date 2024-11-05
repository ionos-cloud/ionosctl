package gateway

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard gateway",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List WireGuard Gateways",
		Example:   "ionosctl vpn wireguard gateway list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, err := Gateways(
				func(req vpn.ApiWireguardgatewaysGetRequest) (vpn.ApiWireguardgatewaysGetRequest, error) {
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

			table, err := resource2table.ConvertVPNWireguardGatewaysToTable(ls)
			if err != nil {
				return fmt.Errorf("could not convert from JSON to Table format: %w", err)
			}
			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutputPreconverted(ls, table,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
	})

	return cmd
}
