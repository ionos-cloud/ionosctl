package file

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "file",
			Aliases:          []string{"f"},
			Short:            "All commands related to zone files",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(getCmd())
	cmd.AddCommand(updateCmd())

	return cmd
}
