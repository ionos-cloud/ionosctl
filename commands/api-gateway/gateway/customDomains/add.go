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
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
	"github.com/spf13/viper"
)

func AddCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigateway",
		Resource:  "customdomains",
		Verb:      "add",
		Example:   "ionosctl apigateway gateway customdomains add --gateway-id ID --name NAME --certificate-id CERTIFICATEID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagName, constants.FlagCertificateId); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			apigatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			usedGateway, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysFindById(context.Background(), apigatewayId).Execute()
			input := usedGateway.Properties
			elem := len(input.CustomDomains)

			if input.CustomDomains == nil || elem == 0 {
				input.CustomDomains = make([]apigateway.GatewayCustomDomains, 1)
				elem = 0
			} else {
				input.CustomDomains = append(input.CustomDomains, apigateway.GatewayCustomDomains{})
				elem = len(input.CustomDomains) - 1
			}

			// max 5 elements
			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.CustomDomains[elem].Name = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagCertificateId); viper.IsSet(fn) {
				input.CustomDomains[elem].CertificateId = pointer.From(viper.GetString(fn))
			}

			input.Name = usedGateway.Properties.Name
			input.Metrics = usedGateway.Properties.Metrics
			input.Logs = usedGateway.Properties.Logs

			gat, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysPut(context.Background(), apigatewayId).
				GatewayEnsure(apigateway.GatewayEnsure{
					Id:         apigatewayId,
					Properties: input,
				}).Execute()

			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			customDomainsConverted := resource2table.ConvertApiGatewayCustomDomainsToTable(gat.Properties.CustomDomains)

			out, err := jsontabwriter.GenerateOutputPreconverted(gat.Properties.CustomDomains, customDomainsConverted,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})
	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return ApiGatewayRouteCreateFlags(cmd)
}
func ApiGatewayRouteCreateFlags(cmd *core.Command) *core.Command {
	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagGatewayShort, "", constants.DescGateway, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.GatewaysIDs()
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)
	cmd.AddStringFlag(constants.FlagName, "", "", "The domain name of the distribution. Field is validated as FQDN according to RFC1123.", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagCertificateId, "", "", "The ID of the certificate to use for the distribution.", core.RequiredFlagOption())
	return cmd
}
