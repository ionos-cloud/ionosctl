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

The domain model has three resources, each with its own sub-command:

  * group  - the scaling group. It owns three things: (1) the replica count bounds (minReplicaCount/maxReplicaCount), (2) a scaling policy that watches a metric (CPU or network) and defines when/how to scale in and out, and (3) a replicaConfiguration template (cores, RAM, CPU family, NICs, volumes) that every new replica is cloned from. A group lives in one datacenter/location.
  * server - the individual VMs currently running inside a group. These are created and destroyed by the autoscaler, not by you directly. ionosctl enriches each entry with its underlying CloudAPI server details.
  * action - the audit log of scaling events (SCALE_IN / SCALE_OUT), each with a status (IN_PROGRESS / SUCCESSFUL / FAILED). Use this to see when and why the group last resized.

A typical flow: create a group (with policy + replica template), let it scale automatically, then inspect 'server' to see live replicas and 'action' to see the scaling history.`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(group.Root())
	cmd.AddCommand(server.Root())
	cmd.AddCommand(action.Root())

	return core.WithConfigOverride(cmd, []string{fileconfiguration.Autoscaling, "vmautoscaling"}, constants.DefaultApiURL+"/autoscaling")
}
