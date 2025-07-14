package vulnerabilities

import (
	"context"
	"fmt"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/artifacts"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VulnerabilitiesListCmd() *core.Command {
	c := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "vulnerabilities",
			Verb:       "list",
			ShortDesc:  "Retrieve vulnerabilities",
			LongDesc:   "Retrieve all vulnerabilities from an artifact",
			Example:    "ionosctl container-registry vulnerabilities list",
			PreCmdRun:  PreCmdList,
			CmdRun:     CmdList,
			InitClient: true,
		},
	)

	c.Command.Flags().StringSlice(constants.FlagCols, nil, tabheaders.ColsMessage(allCols))
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagCols,
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

	c.AddStringFlag("repository", "", "", "Name of the repository to retrieve artifact from")
	_ = c.Command.RegisterFlagCompletionFunc(
		"repository", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return repository.RepositoryNames(viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))),
				cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagArtifactId, "", "", "ID/digest of the artifact")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagArtifactId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return artifacts.ArtifactsIds(
					viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId)),
					viper.GetString(core.GetFlagName(c.NS, "repository")),
				),
				cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddSetFlag(
		cloudapiv6.FlagOrderBy, "", "-score", []string{
			"-score", "-severity", "-publishedAt", "-updatedAt", "-fixable", "score",
			"severity", "publishedAt", "updatedAt", "fixable",
		}, cloudapiv6.FlagOrderByDescription,
	)
	c.AddStringSliceFlag(
		cloudapiv6.FlagFilters, cloudapiv6.FlagFiltersShort, []string{""}, cloudapiv6.FlagFiltersDescription,
	)
	c.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, 100, "Maximum number of results to display")

	return c
}

func PreCmdList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagRegistryId, "repository", constants.FlagArtifactId,
	); err != nil {
		return err
	}

	return query.ValidateFilters(c, []string{"severity", "fixable"}, "Filters available: severity, fixable")
}

func CmdList(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)
	registryId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))
	repository := viper.GetString(core.GetFlagName(c.NS, "repository"))
	artifactId := viper.GetString(core.GetFlagName(c.NS, constants.FlagArtifactId))

	queryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	vulnerabilities, _, err := buildListRequest(registryId, repository, artifactId, queryParams).Execute()
	if err != nil {
		return err
	}

	vulnerabilitiesConverted, err := resource2table.ConvertContainerRegistryVulnerabilitiesToTable(
		vulnerabilities,
	)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(
		vulnerabilities, vulnerabilitiesConverted,
		tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func buildListRequest(
	registryId string, repository string, artifactId string, queryParams resources.ListQueryParams,
) containerregistry.ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest {
	if structs.IsZero(queryParams) {
		return client.Must().RegistryClient.ArtifactsApi.
			RegistriesRepositoriesArtifactsVulnerabilitiesGet(
				context.Background(), registryId, repository, artifactId,
			)
	}

	req := client.Must().RegistryClient.ArtifactsApi.RegistriesRepositoriesArtifactsVulnerabilitiesGet(
		context.Background(), registryId, repository, artifactId,
	)

	if queryParams.OrderBy != nil {
		req = req.OrderBy(*queryParams.OrderBy)
	}

	if queryParams.MaxResults != nil {
		req = req.Limit(*queryParams.MaxResults)
	}

	if queryParams.Filters != nil {
		severity, ok := (*queryParams.Filters)["severity"]
		if ok {
			req = req.FilterSeverity(severity[0])
		}

		fixable, ok := (*queryParams.Filters)["fixable"]
		if ok && fixable[0] == "true" {
			req = req.FilterFixable(true)
		} else if ok && fixable[0] == "false" {
			req = req.FilterFixable(false)
		}
	}

	return req
}
