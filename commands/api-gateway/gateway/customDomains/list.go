package customDomains

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func ListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigateway",
		Resource:  "customdomains",
		Verb:      "list",
		Example:   "ionosctl apigateway gateway customdomains list --gateway-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			apiGatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			rec, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysFindById(context.Background(), apiGatewayId).Execute()
			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

			customDomainsConverted := resource2table.ConvertApiGatewayCustomDomainsToTable(rec.Properties.CustomDomains)

			out, err := jsontabwriter.GenerateOutputPreconverted(rec.Properties.CustomDomains, customDomainsConverted,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagGatewayShort, "", constants.DescGateway, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.GatewaysIDs()
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	return cmd
}
