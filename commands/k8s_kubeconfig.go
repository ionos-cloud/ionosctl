package commands

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func k8sKubeconfig() *builder.Command {
	ctx := context.TODO()
	k8sCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "k8s-kubeconfig",
			Short:            "K8s Kubeconfig Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl k8s-kubeconfig` + "`" + ` allows you to get the Kubeconfig file of a Cluster.`,
			TraverseChildren: true,
		},
	}
	get := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterIdValidate, RunK8sKubeconfigGet, "get", "Get a K8s Cluster Kubeconfig",
		"Use this command to retrieve the kubeconfig file for a given Kubernetes Cluster.\n\nRequired values to run command:\n\n* K8s Cluster Id",
		getK8sKubeconfigExample, true)
	get.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return k8sCmd
}

func RunK8sKubeconfigGet(c *builder.CommandConfig) error {
	u, _, err := c.K8s().ReadKubeConfig(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getKubeFile(u))
}

func getKubeFile(u resources.K8sKubeconfig) string {
	if properties, ok := u.GetPropertiesOk(); ok && properties != nil {
		if kubefile, ok := properties.GetKubeconfigOk(); ok && kubefile != nil {
			return *kubefile
		}
	}
	return ""
}
