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
)

func GatewayPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigateway",
		Resource:  "gateway",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a gateway's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		Example:   "ionosctl apigateway gateway update --gateway-id GATEWAYID --logs true",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			apigatewayId, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
			if err != nil {
				return err
			}

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
	cmd.AddStringFlag(constants.FlagNameCustomDomainsName, "", "", "The domain name of the distribution. Field is validated as FQDN according to RFC1123.")
	cmd.AddStringFlag(constants.FlagCustomCertificateId, "", "", "The ID of the certificate to use for the distribution.")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func partiallyUpdateGatewayPrint(c *core.CommandConfig, r apigateway.GatewayRead) error {
	input := r.Properties
	if c.Command.Command.Flags().Changed(constants.FlagName) {
		name, err := c.Command.Command.Flags().GetString(constants.FlagName)
		if err != nil {
			return err
		}
		input.Name = name
	}
	if c.Command.Command.Flags().Changed(constants.FlagLogs) {
		logs, err := c.Command.Command.Flags().GetBool(constants.FlagLogs)
		if err != nil {
			return err
		}
		input.Logs = pointer.From(logs)
	}
	if c.Command.Command.Flags().Changed(constants.FlagMetrics) {
		metrics, err := c.Command.Command.Flags().GetBool(constants.FlagMetrics)
		if err != nil {
			return err
		}
		input.Metrics = pointer.From(metrics)
	}
	if c.Command.Command.Flags().Changed(constants.FlagNameCustomDomainsName) {
		if len(input.CustomDomains) == 0 {
			input.CustomDomains = make([]apigateway.GatewayCustomDomains, 1)
		}
		name, err := c.Command.Command.Flags().GetString(constants.FlagNameCustomDomainsName)
		if err != nil {
			return err
		}
		input.CustomDomains[0].Name = pointer.From(name)
	}
	if c.Command.Command.Flags().Changed(constants.FlagCustomCertificateId) {
		if len(input.CustomDomains) == 0 {
			input.CustomDomains = make([]apigateway.GatewayCustomDomains, 1)
		}
		certID, err := c.Command.Command.Flags().GetString(constants.FlagCustomCertificateId)
		if err != nil {
			return err
		}
		input.CustomDomains[0].CertificateId = pointer.From(certID)
	}
	apigatewayid, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
	if err != nil {
		return err
	}
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}
