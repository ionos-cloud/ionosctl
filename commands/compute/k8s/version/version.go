package version

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func K8sVersionCmd() *core.Command {
	k8sVersionCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "version",
			Aliases:          []string{"v"},
			Short:            "Kubernetes Version Operations",
			Long:             "The sub-commands of `ionosctl k8s version` allow you to get information about available Kubernetes versions.",
			TraverseChildren: true,
		},
	}

	k8sVersionCmd.AddCommand(K8sVersionListCmd())
	k8sVersionCmd.AddCommand(K8sVersionGetCmd())

	return core.WithConfigOverride(k8sVersionCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
