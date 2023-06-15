package cfg

import (
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
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
	cmd.AddCommand(CfgLocationCmd())
	cmd.AddCommand(LoginCmd())
	cmd.AddCommand(LogoutCmd())
	cmd.AddCommand(WhoamiCmd())
	return cmd
}
