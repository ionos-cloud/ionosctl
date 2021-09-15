package cloudapi_v5

import (
	"context"
	"os"

	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func K8sKubeconfigCmd() *core.Command {
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
	get.AddStringFlag(cloudapiv5.ArgK8sClusterId, "", "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return k8sCmd
}

func RunK8sKubeconfigGet(c *core.CommandConfig) error {
	c.Printer.Verbose("K8s kube config with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
	u, resp, err := c.CloudApiV5Services.K8s().ReadKubeConfig(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
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
