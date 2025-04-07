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
			Short:            "An API gateway consists of generic rules and configurations",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(gateway.GatewayCommand())
	cmd.AddCommand(route.RecordCommand())
	//cmd.AddCommand(reverse_record.Root())
	//cmd.AddCommand(quota.Root())
	//cmd.AddCommand(dnssec.Root())
	//cmd.AddCommand(secondary_zones.Root())

	return core.WithRegionalFlags(cmd, constants.ApiGatewayRegionalURL, constants.GatewayLocations)
}
