package commands

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/spf13/cobra"
)

func k8sVersion() *builder.Command {
	ctx := context.TODO()
	k8sCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "k8s-version",
			Short:            "K8s Version Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl k8s-version` + "`" + ` allows you to get the Kubeconfig file of a Cluster.`,
			TraverseChildren: true,
		},
	}

	builder.NewCommand(ctx, k8sCmd, noPreRun, RunK8sVersionList, "list", "List K8s Versions",
		"Use this command to retrieve all available Kubernetes versions.", listK8sVersionsExample, true)

	builder.NewCommand(ctx, k8sCmd, noPreRun, RunK8sVersionGet, "get", "Get K8s Default Version",
		"Use this command to retrieve the current default Kubernetes version for Clusters and NodePools.",
		getK8sVersionExample, true)

	return k8sCmd
}

func RunK8sVersionList(c *builder.CommandConfig) error {
	u, _, err := c.K8s().ListVersions()
	if err != nil {
		return err
	}
	return c.Printer.Print(u)
}

func RunK8sVersionGet(c *builder.CommandConfig) error {
	u, _, err := c.K8s().GetVersion()
	if err != nil {
		return err
	}
	return c.Printer.Print(u)
}
