package api_gateway

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/gateway"
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/route"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "apigateway",
			Short:            "API Gateway is a service that allows you to monitor API usage, track performance metrics, and generate logs for analysis and troubleshooting.",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(gateway.GatewayCommand())
	cmd.AddCommand(route.RecordCommand())

	return core.WithRegionalFlags(cmd, "apigateway", constants.ApiGatewayRegionalURL, constants.GatewayLocations)
}
