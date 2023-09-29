package commands

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
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
	get.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)
	return k8sCmd
}

func RunK8sKubeconfigGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"K8s kube config with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))

	u, resp, err := c.CloudApiV6Services.K8s().ReadKubeConfig(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateRawOutput(u))

	return nil
}
