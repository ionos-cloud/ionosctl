package cfg

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func ConfigCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "config",
			Short:            "Manage your CLI Configuration",
			Aliases:          []string{"cfg"},
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(LocationCmd())
	cmd.AddCommand(Login())
	cmd.AddCommand(LogoutCmd())
	cmd.AddCommand(WhoamiCmd())

	// adding --api-url support in other ways than 'WithConfigOverride'
	// adds a few edge cases which are a pain to treat, and this is convenient
	return cmd
}
