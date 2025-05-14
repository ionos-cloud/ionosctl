package v1

import (
	cfgv2 "github.com/ionos-cloud/ionosctl/v6/commands/cfg/v2"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "config-v1",
			Short:            "DEPRECATED - Manage your CLI Configuration",
			Aliases:          []string{"cfg-v1", "cfgv1"},
			TraverseChildren: true,
		},
	}

	cmd.Command.Deprecated = "Use 'cfg' instead"

	cmd.AddCommand(cfgv2.LocationCmd())
	cmd.AddCommand(LoginCmd())
	cmd.AddCommand(LogoutCmd())
	cmd.AddCommand(cfgv2.WhoamiCmd())
	return cmd
}
