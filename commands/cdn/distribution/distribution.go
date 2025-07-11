package distribution

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/distribution/routingrules"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols     = []string{"Id", "Domain", "CertificateId", "State"}
	defaultCols = []string{"Id", "Domain", "CertificateId", "State"}
)

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "distribution",
			Short:            "The sub-commands of 'ionosctl cdn distribution' allow you to manage CDN distributions",
			Aliases:          []string{"ds"},
			TraverseChildren: true,
		},
	}
	cmd.Command.PersistentFlags().StringSlice(constants.FlagCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(List())
	cmd.AddCommand(FindByID())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Create())
	cmd.AddCommand(Update())
	cmd.AddCommand(routingrules.Root())
	return cmd
}
