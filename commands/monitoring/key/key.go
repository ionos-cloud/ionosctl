package key

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func KeyCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "key",
			Aliases: []string{"k"},
			Short:   "Rotate a pipeline's ingest key",
			Long: `Manage the ingest key (API key) of a monitoring pipeline. This key authenticates metric pushes: agents send it when writing to the pipeline's HTTP endpoint.

The full key is returned only once - at pipeline creation, and again each time it is rotated here. Rotating generates a fresh key and immediately invalidates the previous one, so any agent still using the old key will start failing until reconfigured.`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(KeyPostCmd())

	return cmd
}
