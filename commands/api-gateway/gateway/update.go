package gateway

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func GatewayPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "api-gateway",
		Resource:  "gateway",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a gateway's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		Example:   "ionosctl apigateway gateway update --gateway-id GATEWAYID --name NAME",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			apigatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			g, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysFindById(context.Background(), apigatewayId).Execute()
			if err != nil {
				return err
			}
			return partiallyUpdateGatewayPrint(c, g)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagGatewayShort, "", constants.DescGateway, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.GatewaysIDs()
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The new name of the ApiGateway", core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.FlagLogs, "", false, "This field enables or disables the collection and reporting of logs for observability of this instance.")
	cmd.AddBoolFlag(constants.FlagMetrics, "", false, "This field enables or disables the collection and reporting of metrics for observability of this instance.")
	cmd.AddStringFlag(constants.FlagFilterName, "", "", "The domain name of the distribution. Field is validated as FQDN according to RFC1123.")
	cmd.AddStringFlag(constants.FlagCertificateId, "", "", "The ID of the certificate to use for the distribution.")

	//cmd.Command.SilenceUsage = true
	//cmd.Command.Flags().SortFlags = false

	return cmd
}

func partiallyUpdateGatewayPrint(c *core.CommandConfig, r apigateway.GatewayRead) error {
	input := r.Properties
	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		input.Name = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagLogs); true {
		input.Logs = pointer.From(viper.GetBool(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagMetrics); true {
		input.Metrics = pointer.From(viper.GetBool(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagFilterName); viper.IsSet(fn) {
		input.CustomDomains[0].Name = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagCertificateId); viper.IsSet(fn) {
		input.CustomDomains[0].CertificateId = pointer.From(viper.GetString(fn))
	}
	apigatewayid := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
	rn, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysPut(context.Background(), apigatewayid).
		GatewayEnsure(apigateway.GatewayEnsure{
			Id:         apigatewayid,
			Properties: input,
		}).Execute()

	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.ApiGatewayGateway, rn,
		tabheaders.GetHeadersAllDefault(allCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
	// TODO
}
