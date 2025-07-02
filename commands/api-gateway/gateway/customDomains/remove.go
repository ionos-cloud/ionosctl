package customDomains

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
	"github.com/spf13/viper"
)

func RemovetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigateway",
		Resource:  "customdomains",
		Verb:      "remove",
		Example:   "ionosctl apigateway gateway customdomains remove --gateway-id ID --custom-domains-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagCustomDomainsId); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			apigatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			customDomainsId := viper.GetInt(core.GetFlagName(c.NS, constants.FlagCustomDomainsId))
			usedApiGateway, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysFindById(context.Background(), apigatewayId).Execute()
			if err != nil {
				return err
			}
			input := usedApiGateway.Properties
			if input.CustomDomains == nil || len(input.CustomDomains) == 0 {
				return fmt.Errorf("there are no custom domains defined in this API Gateway")
			}

			if customDomainsId < 0 || customDomainsId >= len(input.CustomDomains) {
				return fmt.Errorf("invalid custom domain index")
			}

			input.CustomDomains = append(input.CustomDomains[:customDomainsId], input.CustomDomains[customDomainsId+1:]...)
			_, _, err = client.Must().Apigateway.APIGatewaysApi.ApigatewaysPut(context.Background(), apigatewayId).
				GatewayEnsure(apigateway.GatewayEnsure{
					Id:         apigatewayId,
					Properties: input,
				}).Execute()
			if err != nil {
				return err
			}
			// the maximum number of custom domains is 5 (allowed by API)
			return nil
		},
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagGatewayShort, "", constants.DescGateway, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.GatewaysIDs()
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)
	cmd.AddStringFlag(constants.FlagCustomDomainsId, "", "", "The ID of the custom domain", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			apigateway := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID))
			return completer.CustomDomainsIDs(apigateway)
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations))

	return cmd
}
