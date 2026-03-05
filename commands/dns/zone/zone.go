package zone

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone/file"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.zoneName", Default: true},
	{Name: "Description", JSONPath: "properties.description", Default: true},
	{Name: "NameServers", JSONPath: "metadata.nameServers", Default: true},
	{Name: "Enabled", JSONPath: "properties.enabled", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func ZoneCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "zone",
			Aliases:          []string{"z", "zones"},
			Short:            "The sub-commands of 'ionosctl dns zone' allow you to manage DNS zones. A DNS zone serves as an authoritative source of information about which IP addresses belong to which domains",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(ZonesGetCmd())
	cmd.AddCommand(ZonesDeleteCmd())
	cmd.AddCommand(ZonesPostCmd())
	cmd.AddCommand(ZonesPutCmd())
	cmd.AddCommand(ZonesFindByIdCmd())
	cmd.AddCommand(file.Root())

	return cmd
}
