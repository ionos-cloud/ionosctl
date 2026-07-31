package kubeconfig

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func K8sKubeconfigCmd() *core.Command {
	k8sKubeconfigCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "kubeconfig",
			Aliases: []string{"cfg", "config"},
			Short:   "Kubernetes Kubeconfig Operations",
			Long: `Retrieve the kubeconfig file for a Managed Kubernetes cluster.

The kubeconfig contains the API-server endpoint and the credentials kubectl (and
other Kubernetes tooling) need to talk to the cluster. Fetch it once the cluster
is ACTIVE and point kubectl at it, e.g. by saving it to a file and exporting
KUBECONFIG.`,
			TraverseChildren: true,
		},
	}

	k8sKubeconfigCmd.AddCommand(K8sKubeconfigGetCmd())

	return core.WithConfigOverride(k8sKubeconfigCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
