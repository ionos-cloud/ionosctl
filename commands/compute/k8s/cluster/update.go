package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func K8sClusterUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Kubernetes Cluster (control plane)",
		LongDesc: `Update an existing cluster's control-plane settings: name, Kubernetes version,
maintenance window (day + time), API subnet allow list, and S3 audit log bucket.

Upgrading the Kubernetes version: --k8s-version may only be raised, and only to
a version the API currently offers as an upgrade target for this cluster (tab
completion lists them). Downgrades are not supported. Upgrade the cluster before
upgrading its node pools, since a node pool may never exceed the cluster version.

Maintenance window: --maintenance-day and --maintenance-time together define a
weekly window during which IONOS may apply patches and minor upgrades. Setting a
new window without both parts leaves the missing part unchanged.

Note: --api-subnets and --s3bucket overwrite the previous values rather than
merging with them.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the cluster returns to the AVAILABLE state.

Required values to run command:

* K8s Cluster Id`,
		Example: `# Rename a cluster
ionosctl compute k8s cluster update --cluster-id CLUSTER_ID --name new-name

# Upgrade the control-plane Kubernetes version and set a Sunday-night maintenance window
ionosctl compute k8s cluster update --cluster-id CLUSTER_ID --k8s-version 1.29.5 --maintenance-day Sunday --maintenance-time 02:00:00`,
		PreCmdRun:  PreRunK8sClusterId,
		CmdRun:     RunK8sClusterUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "New name for the cluster. Up to 63 characters; must start and end with an alphanumeric character")
	cmd.AddStringFlag(cloudapiv6.ArgK8sVersion, "", "", "New Kubernetes control-plane version, e.g. 1.29.5. Upgrade-only (no downgrades) and must be an offered upgrade target (see tab completion). Upgrade the cluster before its node pools")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sVersion,
		func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
			clusterId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))
			return completer.K8sClusterUpgradeVersions(clusterId), cobra.ShellCompDirectiveNoFileComp
		})
	cmd.AddStringFlag(cloudapiv6.ArgS3Bucket, "", "", "Name of the IONOS Object Storage (S3) bucket that receives API-server audit logs. Overwrites the previous value")
	cmd.AddStringSliceFlag(cloudapiv6.ArgApiSubnets, "", []string{""}, "Restrict access to the Kubernetes API server to these CIDRs (allow list). Cluster-internal traffic is not affected. If empty, access is not restricted. An IP without a subnet mask defaults to /32 for IPv4 and /128 for IPv6. Overwrites the existing allow list")
	cmd.AddStringFlag(cloudapiv6.ArgK8sMaintenanceDay, "", "", "Day of the week for the maintenance window, in English (e.g. Monday, Saturday). Set together with --maintenance-time to define the weekly window")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sMaintenanceDay, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgK8sMaintenanceTime, "", "", "Start time of the maintenance window in HH:mm:ss (24-hour) format, e.g. 08:00:00. Set together with --maintenance-day")
	cmd.AddUUIDFlag(constants.FlagClusterId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgPublic, "", true, "Whether the cluster is public or private. Note: the public/private nature of a cluster is fixed at creation and cannot be changed afterwards")

	return cmd
}
