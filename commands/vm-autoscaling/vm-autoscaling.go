package vm_autoscaling

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/action"
	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/group"
	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/server"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "vm-autoscaling",
			Aliases: []string{"vmas", "vm-as", "vmasc", "vm-asc", "vmautoscaling"},
			Short:   "Manage VM Auto Scaling groups, their servers, and scaling actions",
			Long: `VM Auto Scaling keeps a fleet of identical VMs (replicas) sized to demand: it adds servers when load rises and removes them when load falls, keeping the replica count between a configured minimum and maximum.

A scaling group bundles the replica count bounds, a scaling policy that watches a metric (CPU or network) to decide when to scale in and out, and a replica template (cores, RAM, CPU family, NICs, volumes) that every new server is cloned from; the autoscaler creates and destroys servers automatically and records each SCALE_IN/SCALE_OUT event as an action.

A typical flow: create a group (with policy + replica template), let it scale automatically, then inspect its servers and scaling actions.`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(group.Root())
	cmd.AddCommand(server.Root())
	cmd.AddCommand(action.Root())

	return core.WithConfigOverride(cmd, []string{fileconfiguration.Autoscaling, "vmautoscaling"}, constants.DefaultApiURL+"/autoscaling")
}
