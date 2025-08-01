package gateway

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func ApigatewayListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigateway",
		Resource:  "gateway",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "Retrieve gateways",
		Example:   "ionosctl apigateway gateway list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			req := client.Must().Apigateway.APIGatewaysApi.ApigatewaysGet(context.Background())

			if fn := core.GetFlagName(c.NS, constants.FlagOrderBy); viper.IsSet(fn) {
				req = req.OrderBy(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
				req = req.Offset(viper.GetInt32(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(fn) {
				req = req.Limit(viper.GetInt32(fn))
			}

			ls, _, err := req.Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("items", jsonpaths.ApiGatewayGateway, ls,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagOrderBy, "", "", "The field to order the results by. If not provided, the results will be ordered by the default field.")
	cmd.AddInt32Flag(constants.FlagMaxResults, "", 0, "Pagination limit")
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "Pagination offset")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
