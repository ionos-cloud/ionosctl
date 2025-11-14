package artifacts

import (
	"context"
	"fmt"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ArtifactsListCmd() *core.Command {
	c := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "artifacts",
			Verb:       "list",
			Aliases:    []string{"l", "ls"},
			ShortDesc:  "List registry or repository artifacts",
			LongDesc:   "List all artifacts in a registry or repository",
			Example:    "ionosctl container-registry artifacts list",
			PreCmdRun:  PreCmdList,
			CmdRun:     CmdList,
			InitClient: true,
		},
	)

	c.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
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

	c.AddStringFlag("repository", "", "", "Name of the repository to list artifacts from")
	_ = c.Command.RegisterFlagCompletionFunc(
		"repository", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return repository.RepositoryNames(viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))),
				cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List all artifacts in the registry")
	c.AddSetFlag(
		cloudapiv6.ArgOrderBy, "", "-pullcount", []string{
			"-pullcount", "-pushcount", "-lastPush",
			"-lastPull", "-lastScan", "-vulnTotalCount", "-vulnFixableCount", "pullCount", "pushCount", "lastPush",
			"lastPull", "lastScan", "vulnTotalCount", "vulnFixableCount",
		}, cloudapiv6.ArgOrderByDescription,
	)
	c.AddStringSliceFlag(
		cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription,
	)

	return c
}

func PreCmdList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(
		c.Command, c.NS, []string{constants.FlagRegistryId, "repository"},
		[]string{constants.FlagRegistryId, constants.ArgAll},
	); err != nil {
		return err
	}

	if !viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) && viper.IsSet(
		core.GetFlagName(
			c.NS, cloudapiv6.ArgFilters,
		),
	) {
		return fmt.Errorf("flag --%s can only be used with --%s", cloudapiv6.ArgFilters, constants.ArgAll)
	}

	return query.ValidateFilters(c, []string{"vulnerabilityId"}, "Filters available: vulnerabilityId")
}

func CmdList(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	regId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))
	defCols := defaultCols

	var arts interface{}
	var err error

	queryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
		arts, _, err = buildListAllRequest(regId).Execute()
		if err != nil {
			return err
		}

		defCols = append(defCols, "Repository")
	} else {
		repo := viper.GetString(core.GetFlagName(c.NS, "repository"))

		arts, _, err = buildListRequest(regId, repo).Execute()
		if err != nil {
			return err
		}
	}

	out, err := jsontabwriter.GenerateOutput(
		"items", jsonpaths.ContainerRegistryArtifact, arts,
		tabheaders.GetHeaders(allCols, defCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func buildListAllRequest(
	regId string, queryParams resources.ListQueryParams,
) containerregistry.ApiRegistriesArtifactsGetRequest {
	if structs.IsZero(queryParams) {
		return client.Must().RegistryClient.ArtifactsApi.RegistriesArtifactsGet(
			context.Background(), regId,
		)
	}

	req := client.Must().RegistryClient.ArtifactsApi.RegistriesArtifactsGet(context.Background(), regId)

	if queryParams.OrderBy != nil {
		req = req.OrderBy(*queryParams.OrderBy)
	}

	if queryParams.Filters != nil {
		vulnId, ok := (*queryParams.Filters)["vulnerabilityId"]
		if ok {
			req = req.FilterVulnerabilityId(vulnId[0])
		}
	}

	return req
}

func buildListRequest(
	regId string, repo string, queryParams resources.ListQueryParams,
) containerregistry.ApiRegistriesRepositoriesArtifactsGetRequest {
	if structs.IsZero(queryParams) {
		return client.Must().RegistryClient.ArtifactsApi.RegistriesRepositoriesArtifactsGet(
			context.Background(), regId, repo,
		)
	}

	req := client.Must().RegistryClient.ArtifactsApi.RegistriesRepositoriesArtifactsGet(
		context.Background(), regId, repo,
	)

	if queryParams.OrderBy != nil {
		req = req.OrderBy(*queryParams.OrderBy)
	}

	return req
}
