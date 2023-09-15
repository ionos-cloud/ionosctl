package repository

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RegRepoDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "repository",
			Verb:      "repository",
			Aliases:   []string{"rd", "del", "repo", "rep-del", "repository-delete"},
			ShortDesc: "Delete all repository contents.",
			LongDesc: "Delete all repository contents. The registry V2 API allows manifests and blobs to be deleted " +
				"individually but it is not possible to remove an entire repository. This operation is provided for " +
				"convenience",
			Example:    "ionosctl container-registry repository-delete --registry-id [REGISTRY-ID], --name [REPOSITORY-NAME]",
			PreCmdRun:  PreCmdDelete,
			CmdRun:     CmdDelete,
			InitClient: true,
		},
	)

	cmd.AddStringFlag("name", "n", "", "Name of the repository to delete")
	cmd.AddStringFlag("registry-id", "r", "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdDelete(c *core.CommandConfig) error {
	regId := viper.GetString("registry-id")
	repoName := viper.GetString("name")

	res, _ := c.ContainerRegistryServices.Repository().Delete(regId, repoName)
	if res.StatusCode == 204 {
		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Repository is being deleted"))
	} else {
		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Repository could not be deleted"))
	}

	return nil
}

func PreCmdDelete(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, "name", "registry-id")
	if err != nil {
		return err
	}

	return nil
}
