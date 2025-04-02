package gateway

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
	"github.com/spf13/viper"

	//"fmt"
	//"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	//"github.com/spf13/viper"
)

func ApigatewayPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "api-gateway",
		Resource:  "apigateway",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Create an apigateway",
		Example:   "ionosctl apigateway create --name name.com",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {

			//request_smth := client.Must().Apigateway.APIGatewaysApi.ApigatewaysPost(context.Background()).GatewayCreate(
			//	apigateway.GatewayCreate{
			//		Properties: apigateway.Gateway{
			//			Name:          "Example",
			//			Logs:          nil,
			//			Metrics:       nil,
			//			CustomDomains: nil,
			//		},
			//	},
			//)

			//execute, _, err := request_smth.Execute()
			//if err != nil {
			//	return err
			//}
			//j, _ := json.MarshalIndent(execute, "", "  ")
			//fmt.Println(string(j))

			input := apigateway.Gateway{}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Name = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagLogs); viper.IsSet(fn) {
				input.Logs = pointer.From(viper.GetBool(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagMetrics); viper.IsSet(fn) {
				input.Metrics = pointer.From(viper.GetBool(fn))
			}

			z, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysPost(context.Background()).
				GatewayCreate(apigateway.GatewayCreate{Properties: input}).Execute()

			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.ApiGatewayGateway, z, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the ApiGateway gateway")
	cmd.AddBoolFlag(constants.FlagLogs, "", false, "The logs parameter of the ApiGateway gateway")
	cmd.AddBoolFlag(constants.FlagMetrics, "", false, "Activate or deactivate the ApiGateway gateway metrics parameter")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
