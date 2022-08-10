package dataplatform

import (
	"context"
	"os"

	"github.com/ionos-cloud/ionosctl/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	dp "github.com/ionos-cloud/ionosctl/services/dataplatform"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func KubeConfigCmd() *core.Command {
	ctx := context.TODO()
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "kubeconfig",
			Aliases:          []string{"cfg", "config"},
			Short:            "Data Platform Kubeconfig Operations",
			Long:             "The sub-command of `ionosctl dataplatform kubeconfig` allows you to get the configuration file of a Data Platform Cluster.",
			TraverseChildren: true,
		},
	}

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "dataplatform",
		Resource:   "kubeconfig",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get the kubeconfig file for a Data Platform Cluster",
		LongDesc:   "Use this command to retrieve the kubeconfig file for a given Data Platform Cluster.\n\nRequired values to run command:\n\n* K8s Cluster Id",
		Example:    getKubeConfigExample,
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunKubeConfigGet,
		InitClient: true,
	})
	get.AddStringFlag(dp.ArgClusterId, "", "", dp.ClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(dp.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

	return k8sCmd
}

func RunKubeConfigGet(c *core.CommandConfig) error {
	c.Printer.Verbose("K8s kube config with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)))
	u, _, err := c.DataPlatformServices.Clusters().GetKubeConfig(viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(u)
}
