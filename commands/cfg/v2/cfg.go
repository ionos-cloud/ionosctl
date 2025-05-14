package v2

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "config",
			Short:            "Manage your CLI Configuration",
			Aliases:          []string{"cfg"},
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(LocationCmd())
	// cmd.AddCommand(LoginCmd())
	// cmd.AddCommand(LogoutCmd())
	cmd.AddCommand(WhoamiCmd())
	return cmd
}
