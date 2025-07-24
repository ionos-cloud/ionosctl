package key

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func KeyCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "key",
			Aliases:          []string{"k"},
			Short:            "Generates a new key for a pipeline invalidating the old one. The key is used for authentication when sending metrics.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(KeyPostCmd())

	return cmd
}
