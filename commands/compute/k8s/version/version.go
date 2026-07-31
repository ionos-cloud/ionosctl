package version

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func K8sVersionCmd() *core.Command {
	k8sVersionCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "version",
			Aliases: []string{"v"},
			Short:   "Kubernetes Version Operations",
			Long: `Inspect the Kubernetes versions offered by Managed Kubernetes.

These versions are the valid values for --k8s-version when creating or upgrading
clusters and node pools. Remember the skew rule: a node pool's version must be
<= its cluster's version. Use ` + "`list`" + ` to see all offered versions and ` + "`get`" + ` to
see the default that is applied when --k8s-version is omitted.`,
			TraverseChildren: true,
		},
	}

	k8sVersionCmd.AddCommand(K8sVersionListCmd())
	k8sVersionCmd.AddCommand(K8sVersionGetCmd())

	return core.WithConfigOverride(k8sVersionCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
