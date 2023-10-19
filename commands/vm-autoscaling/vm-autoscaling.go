package vm_autoscaling

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/groups"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
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

	cmd.AddCommand(groups.GroupsCmd())

	return cmd
}
