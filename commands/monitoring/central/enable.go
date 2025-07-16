package central

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/central/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func CentralEnable() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "central",
		Verb:      "enable",
		Aliases:   []string{"e"},
		ShortDesc: "enable a CentralMonitoring",
		Example:   "ionosctl monitoring central enable --location de/txl",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			completer.CentralEnable(c, true)
			return nil
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
