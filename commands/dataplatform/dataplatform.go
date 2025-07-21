package dataplatform

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/nodepool"
	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/version"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func DataplatformCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dataplatform",
			Short:            "Managed Stackable Data Platform by IONOS Cloud provides a preconfigured Kubernetes cluster with pre-installed and managed Stackable operators. After the provision of these Stackable operators, the customer can interact with them directly and build his desired application on top of the Stackable Platform.",
			Aliases:          []string{"mdp", "dp", "stackable", "managed-dataplatform"},
			Long:             "The sub-commands of `ionosctl dataplatform` allow you to perform operations on Dataplatform resources.",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(cluster.ClusterCmd())
	cmd.AddCommand(nodepool.NodepoolCmd())
	cmd.AddCommand(version.Root())
	return core.WithConfigOverride(cmd, "dataplatform", constants.DefaultApiURL+"/dataplatform")
}
