package lan

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var (
	defaultNatGatewayLanCols = []string{"NatGatewayLanId", "GatewayIps"}
)

func NatgatewayLanCmd() *core.Command {
	natgatewayLanCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "lan",
			Short:            "NAT Gateway Lan Operations",
			Long:             "The sub-commands of `ionosctl natgateway lan` allow you to add, list, remove NAT Gateway Lans.",
			TraverseChildren: true,
		},
	}

	natgatewayLanCmd.AddCommand(NatgatewayLanListCmd())
	natgatewayLanCmd.AddCommand(NatgatewayLanAddCmd())
	natgatewayLanCmd.AddCommand(NatgatewayLanRemoveCmd())

	return core.WithConfigOverride(natgatewayLanCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
