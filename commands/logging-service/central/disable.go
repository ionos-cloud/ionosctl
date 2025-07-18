package central

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func CentralDisable() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "logging-service",
		Resource:  "central",
		Verb:      "disable",
		Aliases:   []string{"d"},
		ShortDesc: "Disable CentralLogging",
		Example:   "ionosctl logging-service central disable --location de/txl",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
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
