package kubeconfig

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func K8sKubeconfigCmd() *core.Command {
	k8sKubeconfigCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "kubeconfig",
			Aliases:          []string{"cfg", "config"},
			Short:            "Kubernetes Kubeconfig Operations",
			Long:             "The sub-command of `ionosctl k8s kubeconfig` allows you to get the configuration file of a Kubernetes Cluster.",
			TraverseChildren: true,
		},
	}

	k8sKubeconfigCmd.AddCommand(K8sKubeconfigGetCmd())

	return core.WithConfigOverride(k8sKubeconfigCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
