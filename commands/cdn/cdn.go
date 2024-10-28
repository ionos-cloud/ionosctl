package cdn

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/distribution"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cdn",
			Short:            "The sub-commands of the 'cdn' resource help manage CDN distributions",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(distribution.Command())
	return cmd
}