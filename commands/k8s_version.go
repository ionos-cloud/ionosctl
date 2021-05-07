package commands

import (
	"context"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/spf13/cobra"
)

func k8sVersion() *builder.Command {
	ctx := context.TODO()
	k8sCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "version",
			Short:            "Kubernetes Version Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl k8s version` + "`" + ` allow you to get information about available Kubernetes versions.`,
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	builder.NewCommand(ctx, k8sCmd, noPreRun, RunK8sVersionList, "list", "List Kubernetes Versions",
		"Use this command to retrieve all available Kubernetes versions.", listK8sVersionsExample, true)

	/*
		Get Command
	*/
	builder.NewCommand(ctx, k8sCmd, noPreRun, RunK8sVersionGet, "get", "Get Kubernetes Default Version",
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
	u, err := getK8sVersion(c)
	if err != nil {
		return err
	}
	return c.Printer.Print(u)
}

func getK8sVersion(c *builder.CommandConfig) (string, error) {
	if k8sversion, _, err := c.K8s().GetVersion(); err == nil {
		k8sversion = strings.ReplaceAll(k8sversion, "\"", "")
		k8sversion = strings.ReplaceAll(k8sversion, "\n", "")
		return k8sversion, nil
	} else {
		return "", err
	}
}
