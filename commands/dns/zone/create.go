package zone

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "create",
		Aliases:   []string{},
		ShortDesc: "Create a zone",
		Example:   "ionosctl dns zonecreate ",
		PreCmdRun: func(c *core.PreCommandConfig) error {

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			// Implement the actual command logic here
			return nil
		},
		InitClient: true,
	})

	return cmd
}
