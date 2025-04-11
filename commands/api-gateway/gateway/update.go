package gateway

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
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
		Example:   "ionosctl apigateway gateway update --gateway-id GATEWAYID",
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

	cmd.AddStringFlag(constants.FlagGatewayRouteID, "", "", fmt.Sprintf("%s. Required or -%s", constants.DescRoute, constants.ArgAllShort),
		core.WithCompletion(func() []string {
			return completer.GatewaysIDs()
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The new name of the DNS zone, e.g. foo.com")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "The new description of the DNS zone")
	cmd.AddBoolFlag(constants.FlagEnabled, "", true, "Activate or deactivate the DNS zone")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func partiallyUpdateGatewayPrint(c *core.CommandConfig, r apigateway.GatewayRead) error {
	input := r.Properties
	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		input.Name = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagLogs); viper.IsSet(fn) {
		input.Logs = pointer.From(viper.GetBool(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagMetrics); viper.IsSet(fn) {
		input.Metrics = pointer.From(viper.GetBool(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagFilterName); viper.IsSet(fn) {
		input.CustomDomains[0].Name = pointer.From(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagMetrics); viper.IsSet(fn) {
		input.CustomDomains[0].CertificateId = pointer.From(viper.GetString(fn))
	}
	return nil
	// TODO
}
