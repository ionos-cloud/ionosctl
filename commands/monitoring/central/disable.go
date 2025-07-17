package central

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/central/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

//ionosctl monitoring central enable/disable/get
//enable face PUT cu properties: enabled: true
//disable cu false
//si get , face get xD

func CentralDisable() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "central",
		Verb:      "disable",
		Aliases:   []string{"d"},
		ShortDesc: "Disable a CentralMonitoring",
		Example:   "ionosctl monitoring central disable --location de/txl",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			return completer.CentralEnable(c, false)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
