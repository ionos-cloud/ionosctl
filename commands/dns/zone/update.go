package zone

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "update",
		Aliases:   []string{},
		ShortDesc: "Ensure a zone",
		Example:   "ionosctl dns zone update",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			return nil
		},
		InitClient: true,
	})

	return cmd
}
