package central

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func CentralEnable() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "central",
		Verb:      "enable",
		Aliases:   []string{"e"},
		ShortDesc: "Turn on central monitoring for a region",
		LongDesc:  `Enable central monitoring in the region given by --location. Once enabled, other IONOS Cloud products in that region automatically forward their metrics to the Monitoring Service, so you do not have to configure a push agent for each product. Prints the resulting state and Grafana endpoint.`,
		Example:   "ionosctl monitoring central enable --location de/txl",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation()
		},
		CmdRun: func(c *core.CommandConfig) error {
			return enable(c, true)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
