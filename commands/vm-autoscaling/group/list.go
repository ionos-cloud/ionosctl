package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "groups",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List VM Autoscaling Groups. Use a greater '--depth' to see current replica count",
		Example:   "ionosctl vm-autoscaling group list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			ls, err := Groups()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Prefix("items").Print(ls)
		},
	})

	return cmd
}
