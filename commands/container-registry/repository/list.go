package repository

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RepositoryListCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "repository",
			Verb:       "list",
			Aliases:    []string{"ls", "l"},
			ShortDesc:  "Retrieve all repositories.",
			LongDesc:   "Retrieve all repositories in a registry.",
			Example:    "ionosctl container-registry list",
			PreCmdRun:  PreCmdList,
			CmdRun:     CmdList,
			InitClient: true,
		},
	)
	c.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(allCols))
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "Registry ID")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddSetFlag(
		cloudapiv6.ArgOrderBy, "", "-lastPush", []string{
			"-lastPush", "-lastPull", "-artifactCount", "-pullCount", "-pushCount", "name", "lastPush",
			"lastPull", "artifactCount", "pullCount", "pushCount",
		}, cloudapiv6.ArgOrderByDescription,
	)
	c.AddStringSliceFlag(
		cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription,
	)

	return c
}

func PreCmdList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagRegistryId); err != nil {
		return err
	}

	return query.ValidateFilters(
		c, []string{"name", "vulnerabilitySeverity"}, "Filters available: name, vulnerabilitySeverity",
	)
}

func CmdList(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	regId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))

	repos, _, err := client.Must().RegistryClient.RepositoriesApi.RegistriesRepositoriesGet(
		context.Background(), regId).Execute()
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput(
		"items", jsonpaths.ContainerRegistryRepository, repos,
		tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
