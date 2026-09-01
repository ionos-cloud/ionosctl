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

The 'Replicas' column (current number of running VMs) may be empty in the list view, as the API response does not include the group's server entities here. Use 'group get' for the full detail of a single group.`,
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
