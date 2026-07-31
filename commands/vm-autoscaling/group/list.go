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
		ShortDesc: "List all VM Auto Scaling groups in your account",
		LongDesc: `List every VM Auto Scaling group your account can access, one row per group, showing its datacenter, replica bounds, scaling metric and location.

The 'Replicas' column (current number of running VMs) is only populated when the API response includes the group's server entities, which requires a deeper query. Increase the global '--depth' flag (e.g. --depth 2) to have those counts filled in.`,
		Example:   "ionosctl vm-autoscaling group list --depth 2",
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
