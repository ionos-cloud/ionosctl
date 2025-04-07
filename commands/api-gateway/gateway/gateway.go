package gateway

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Name", "Logs", "Metrics", "Enable", "DomainName", "CertificateId", "HttpMethods", "HttpCodes", "Override", "Status", "PublicEndpoint"}
)

func GatewayCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "gateway",
			Aliases:          []string{"a", "api"},
			Short:            "API Gateway is a service that allows you to monitor API usage, track performance metrics, and generate logs for analysis and troubleshooting.",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	//
	//cmd.AddCommand(ZonesGetCmd())
	//cmd.AddCommand(ZonesDeleteCmd())
	cmd.AddCommand(ApigatewayPostCmd())
	cmd.AddCommand(ApigatewayListCmd())
	cmd.AddCommand(ApiGatewayDeleteCmd())
	cmd.AddCommand(GatewaysFindByIdCmd())

	return cmd
}
