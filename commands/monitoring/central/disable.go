package central

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func CentralDisable() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "central",
		Verb:      "disable",
		Aliases:   []string{"d"},
		ShortDesc: "Turn off central monitoring for a region",
		LongDesc:  `Disable central monitoring in the region given by --location. After this, IONOS products stop forwarding their metrics automatically; only metrics you push explicitly with a pipeline's ingest key continue to be collected. Existing pipelines and their data are unaffected. Prints the resulting state.`,
		Example:   "ionosctl monitoring central disable --location de/txl",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation()
		},
		CmdRun: func(c *core.CommandConfig) error {
			return enable(c, false)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
