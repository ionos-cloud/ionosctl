package route

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/route/upstreams"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Name", "Type", "Paths", "Methods", "Host", "Port", "Weight", "Status", "StatusMessage"}
)

func RecordCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "route",
			Short:            "A route is a rule that maps an incoming request to a specific backend service.",
			Aliases:          []string{"r"},
			TraverseChildren: true,
		},
	}
	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddCommand(RouteListCmd())
	cmd.AddCommand(ApiGatewayRouteDeleteCmd())
	cmd.AddCommand(ApiGatewayRoutesPostCmd())
	cmd.AddCommand(RouteFindByIdCmd())
	cmd.AddCommand(RoutesPutCmd())
	cmd.AddCommand(upstreams.UpstreamsCmd())
	return cmd
}
