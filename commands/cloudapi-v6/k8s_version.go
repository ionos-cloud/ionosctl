package commands

import (
	"context"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"strings"

	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/spf13/cobra"
)

func K8sVersionCmd() *core.Command {
	ctx := context.TODO()
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "version",
			Aliases:          []string{"v"},
			Short:            "Kubernetes Version Operations",
			Long:             "The sub-commands of `ionosctl k8s version` allow you to get information about available Kubernetes versions.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "version",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Kubernetes Versions",
		LongDesc:   "Use this command to retrieve all available Kubernetes versions.",
		Example:    listK8sVersionsExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunK8sVersionList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "version",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get Kubernetes Default Version",
		LongDesc:   "Use this command to retrieve the current default Kubernetes version for Clusters and NodePools.",
		Example:    getK8sVersionExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunK8sVersionGet,
		InitClient: true,
	})

	return k8sCmd
}

func RunK8sVersionList(c *core.CommandConfig) error {
	u, resp, err := c.CloudApiV6Services.K8s().ListVersions()
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(u)
}

func RunK8sVersionGet(c *core.CommandConfig) error {
	u, err := getK8sVersion(c)
	if err != nil {
		return err
	}
	return c.Printer.Print(u)
}

func getK8sVersion(c *core.CommandConfig) (string, error) {
	if k8sversion, resp, err := c.CloudApiV6Services.K8s().GetVersion(); err == nil {
		k8sversion = strings.ReplaceAll(k8sversion, "\"", "")
		k8sversion = strings.ReplaceAll(k8sversion, "\n", "")
		if resp != nil {
			c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
		}
		return k8sversion, nil
	} else {
		return "", err
	}
}
