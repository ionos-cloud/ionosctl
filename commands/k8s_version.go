package commands

import (
	"context"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func k8sVersion() *core.Command {
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
		PreCmdRun:  noPreRun,
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
		PreCmdRun:  noPreRun,
		CmdRun:     RunK8sVersionGet,
		InitClient: true,
	})

	return k8sCmd
}

func RunK8sVersionList(c *core.CommandConfig) error {
	u, resp, err := c.K8s().ListVersions()
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
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
	if k8sversion, resp, err := c.K8s().GetVersion(); err == nil {
		c.Printer.Verbose("K8s version with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, config.ArgK8sVersion)))
		k8sversion = strings.ReplaceAll(k8sversion, "\"", "")
		k8sversion = strings.ReplaceAll(k8sversion, "\n", "")
		if resp != nil {
			c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
		}
		return k8sversion, nil
	} else {
		return "", err
	}
}
