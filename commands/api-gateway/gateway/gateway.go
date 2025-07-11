package gateway

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/gateway/customDomains"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Name", "Logs", "Metrics", "Enable", "DomainName", "CertificateId", "HttpMethods", "HttpCodes", "Override", "PublicEndpoint", "Status"}
)

func GatewayCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "gateway",
			Aliases:          []string{"a", "api"},
			Short:            "An API gateway consists of generic rules and configurations",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.FlagCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddCommand(GatewayPutCmd())
	cmd.AddCommand(ApigatewayPostCmd())
	cmd.AddCommand(ApigatewayListCmd())
	cmd.AddCommand(ApiGatewayDeleteCmd())
	cmd.AddCommand(GatewaysFindByIdCmd())
	cmd.AddCommand(customDomains.CustomDomainsCmd())

	return cmd
}
