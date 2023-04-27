package record

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"

	ionoscloud "github.com/ionos-cloud/sdk-go-dnsaas"
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
all, mark individual flags as required

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			// Implement the actual command logic here
		},
		InitClient: true,
	})

	return cmd
}
