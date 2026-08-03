package targetgroup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allTargetGroupTargetCols = []table.Column{
	{Name: "TargetIp", JSONPath: "ip", Default: true},
	{Name: "TargetPort", JSONPath: "port", Default: true},
	{Name: "Weight", JSONPath: "weight", Default: true},
	{Name: "HealthCheckEnabled", JSONPath: "healthCheckEnabled", Default: true},
	{Name: "MaintenanceEnabled", JSONPath: "maintenanceEnabled", Default: true},
}

func TargetGroupTargetCmd() *core.Command {
	targetGroupTargetCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "target",
			Aliases: []string{"t"},
			Short:   "Target Group Target Operations",
			Long: `Manage the targets (backend servers) inside a Target Group.

A target is one backend endpoint identified by an IP address and a port. Each target also carries a --weight (its share of traffic under the group's algorithm), plus per-target --health-check-enabled and --maintenance-enabled switches. The group-level health-check settings define HOW targets are probed; these per-target switches decide WHETHER a given target is probed and whether it currently receives traffic at all.

The API stores targets as an array on the Target Group, so add/remove operations read the current group, modify the target list, and write it back. A target is matched for removal by its IP + port pair.`,
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(context.TODO(), targetGroupTargetCmd, core.CommandBuilder{
		Namespace:  "targetgroup",
		Resource:   "target",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Target Groups Targets",
		LongDesc:   "List the targets (backend servers) belonging to a Target Group, showing each target's IP, port, weight, and its health-check-enabled / maintenance-enabled state.\n\nRequired values to run command:\n\n* Target Group Id",
		Example:    `ionosctl compute targetgroup target list --targetgroup-id TARGET_GROUP_ID`,
		PreCmdRun:  PreRunTargetGroupId,
		CmdRun:     RunTargetGroupTargetList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddColsFlag(allTargetGroupTargetCols)

	/*
		Add Command
	*/
	add := core.NewCommand(context.TODO(), targetGroupTargetCmd, core.CommandBuilder{
		Namespace: "targetgroup",
		Resource:  "target",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Target to a Target Group",
		LongDesc: `Add a backend server (target) to a Target Group. A target is identified by its --ip and --port; the same IP with different ports counts as distinct targets.

--weight controls this target's share of traffic relative to the others. Each target receives load proportional to its weight over the sum of all weights (so weight 2 gets twice the traffic of weight 1). Range is 0-256, default 1. A weight of 0 excludes the target from new load-balancing decisions but still lets it serve existing persistent connections - useful for gracefully draining a server. When sizing by capacity, start in the middle of the range (e.g. 10-100) so you can adjust up or down later.

--health-check-enabled (default true) decides whether this target is probed at all. When off, the target is treated as always available and traffic is sent to it blindly. When on, the target only receives traffic while it passes the group's health check (a connection attempt to the target's own IP and port).

--maintenance-enabled (default false) takes the target out of rotation regardless of health, so no balanced traffic reaches it.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Target Group Id
* Target Ip
* Target Port`,
		Example: `# Add a backend with default weight 1 and health checking on
ionosctl compute targetgroup target add --targetgroup-id TARGET_GROUP_ID --ip 10.0.0.5 --port 8080

# Add a higher-capacity backend that should take twice the traffic
ionosctl compute targetgroup target add --targetgroup-id TARGET_GROUP_ID --ip 10.0.0.6 --port 8080 --weight 2

# Add a backend already in maintenance (registered but not receiving traffic)
ionosctl compute targetgroup target add --targetgroup-id TARGET_GROUP_ID --ip 10.0.0.7 --port 8080 --maintenance-enabled`,
		PreCmdRun:  PreRunTargetGroupIdTargetIpPort,
		CmdRun:     RunTargetGroupTargetAdd,
		InitClient: true,
	})
	add.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddIpFlag(cloudapiv6.ArgIp, "", nil, "The IP address of the backend server that will receive balanced traffic.", core.RequiredFlagOption())
	add.AddIntFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, 8080, "The port on the backend server that receives traffic. Valid range is 1 to 65535. Together with --ip it uniquely identifies the target.", core.RequiredFlagOption())
	add.AddIntFlag(cloudapiv6.ArgWeight, cloudapiv6.ArgWeightShort, 1, "This target's share of traffic relative to the combined weight of all targets (higher weight = larger share). Valid range is 0 to 256, default 1. Weight 0 excludes the target from new balancing decisions but still serves existing persistent connections (useful for draining). Prefer mid-range values (e.g. 10-100) to leave room for later adjustment.")
	add.AddBoolFlag(cloudapiv6.ArgHealthCheckEnabled, "", true, "When true (default), the target only receives traffic while it passes the group's health check (a TCP connection attempt to this target's IP and port). When false, the target is treated as always available and is never probed.")
	add.AddBoolFlag(cloudapiv6.ArgMaintenanceEnabled, cloudapiv6.ArgMaintenanceShort, false, "When true, the target is held out of rotation and receives no balanced traffic regardless of its health status. Default is false.")
	add.AddColsFlag(allTargetGroupTargetCols)

	/*
		Remove Command
	*/
	remove := core.NewCommand(context.TODO(), targetGroupTargetCmd, core.CommandBuilder{
		Namespace: "targetgroup",
		Resource:  "target",
		Verb:      "remove",
		Aliases:   []string{"r"},
		ShortDesc: "Remove a Target from a Target Group",
		LongDesc:  "Remove a target (backend server) from a Target Group. The target is matched by its --ip and --port pair; both must match an existing target or the command reports it was not found. Use --all to remove every target from the group, leaving the group defined but empty.\n\nRequired values to run command:\n\n* Target Group Id\n* Target Ip\n* Target Port",
		Example: `# Remove one backend by IP + port
ionosctl compute targetgroup target remove --targetgroup-id TARGET_GROUP_ID --ip 10.0.0.5 --port 8080

# Empty the group (remove all targets)
ionosctl compute targetgroup target remove --targetgroup-id TARGET_GROUP_ID --all --force`,
		PreCmdRun:  PreRunTargetGroupTargetRemove,
		CmdRun:     RunTargetGroupTargetRemove,
		InitClient: true,
	})
	remove.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = remove.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	remove.AddIpFlag(cloudapiv6.ArgIp, "", nil, "The IP address of the target to remove. Must match an existing target's IP together with --port.", core.RequiredFlagOption())
	remove.AddIntFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, 8080, "The port of the target to remove. Together with --ip it identifies which target is removed. Valid range is 1 to 65535.", core.RequiredFlagOption())
	remove.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all targets from the group, leaving it empty. Cannot be combined with --ip / --port.")
	remove.AddColsFlag(allTargetGroupTargetCols)

	return core.WithConfigOverride(targetGroupTargetCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
