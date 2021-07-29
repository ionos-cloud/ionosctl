package commands

import (
	"context"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func k8sKubeconfig() *core.Command {
	ctx := context.TODO()
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "kubeconfig",
			Aliases:          []string{"cfg", "config"},
			Short:            "Kubernetes Kubeconfig Operations",
			Long:             "The sub-command of `ionosctl k8s kubeconfig` allows you to get the configuration file of a Kubernetes Cluster.",
			TraverseChildren: true,
		},
	}

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "kubeconfig",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get the kubeconfig file for a Kubernetes Cluster",
		LongDesc:   "Use this command to retrieve the kubeconfig file for a given Kubernetes Cluster.\n\nRequired values to run command:\n\n* K8s Cluster Id",
		Example:    getK8sKubeconfigExample,
		PreCmdRun:  PreRunK8sClusterId,
		CmdRun:     RunK8sKubeconfigGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return k8sCmd
}

func RunK8sKubeconfigGet(c *core.CommandConfig) error {
	c.Printer.Infof("K8s kube config with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)))
	u, _, err := c.K8s().ReadKubeConfig(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)))
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
