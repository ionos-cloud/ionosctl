package repository

import (
	"context"
	"fmt"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
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

	c.AddStringFlag("registry-id", "r", "", "Registry ID")
	_ = c.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	c.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, 100, "Maximum number of results to display")

	return c
}

func PreCmdList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, "registry-id"); err != nil {
		return err
	}

	return query.ValidateFilters(
		c, []string{"name", "vulnerabilitySeverity"}, "Filters available: name, vulnerabilitySeverity",
	)
}

func CmdList(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	regId := viper.GetString(core.GetFlagName(c.NS, "registry-id"))

	queryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	repos, _, err := buildListRequest(regId, queryParams).Execute()
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func buildListRequest(
	registryId string, queryParams resources.ListQueryParams,
) ionoscloud.
	ApiRegistriesRepositoriesGetRequest {
	if structs.IsZero(queryParams) {
		return client.Must().RegistryClient.RepositoriesApi.RegistriesRepositoriesGet(
			context.Background(),
			registryId,
		)
	}

	req := client.Must().RegistryClient.RepositoriesApi.RegistriesRepositoriesGet(context.Background(), registryId)

	if queryParams.OrderBy != nil {
		req = req.OrderBy(*queryParams.OrderBy)
	}

	if queryParams.MaxResults != nil {
		req = req.Limit(*queryParams.MaxResults)
	}

	if queryParams.Filters != nil {
		vulnSeverity, ok := (*queryParams.Filters)["vulnerabilitySeverity"]
		if ok {
			req = req.FilterVulnerabilitySeverity(vulnSeverity[0])
		}

		name, ok := (*queryParams.Filters)["name"]
		if ok {
			req = req.FilterName(name[0])
		}
	}

	return req
}
