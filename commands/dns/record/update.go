package record

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesRecordsPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "update",
		Aliases:   []string{},
		ShortDesc: "Modify an existing record",
		Example:   "ionosctl dns record update ",
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