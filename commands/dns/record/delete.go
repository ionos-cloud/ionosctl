package record

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesRecordsDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a record",
		Example:   "ionosctl dns record delete ",
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
