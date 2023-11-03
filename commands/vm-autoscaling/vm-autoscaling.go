package vm_autoscaling

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/action"
	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/group"
	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/server"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "vm-autoscaling",
			Aliases:          []string{"vmasc", "vm-asc", "vmautoscaling"},
			Short:            "VM Autoscaling Operations",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(group.Root())
	cmd.AddCommand(server.Root())
	cmd.AddCommand(action.Root())

	return cmd
}
