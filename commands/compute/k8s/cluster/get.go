package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func K8sClusterGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "cluster",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Kubernetes Cluster",
		LongDesc:   "Use this command to retrieve details about a specific Kubernetes Cluster.You can wait for the Cluster to be in \"ACTIVE\" state using `--wait-for-state` flag together with `--timeout` option.\n\nRequired values to run command:\n\n* K8s Cluster Id",
		Example:    "ionosctl k8s cluster get --cluster-id CLUSTER_ID",
		PreCmdRun:  PreRunK8sClusterId,
		CmdRun:     RunK8sClusterGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(constants.FlagClusterId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for specified Cluster to be in ACTIVE state")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.K8sTimeoutSeconds, "Timeout option for waiting for Cluster to be in ACTIVE state [seconds]")

	return cmd
}
