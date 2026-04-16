package artifacts

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
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
			c.NS, constants.FlagFilters,
		),
	) {
		return fmt.Errorf("flag --%s can only be used with --%s", constants.FlagFilters, constants.ArgAll)
	}

	return nil
}

func CmdList(c *core.CommandConfig) error {
	cols := c.Cols()
	regId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))

	var arts interface{}
	var err error
	addRepository := false

	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
		arts, _, err = client.Must().RegistryClient.ArtifactsApi.RegistriesArtifactsGet(
			context.Background(), regId).Execute()
		if err != nil {
			return err
		}

		addRepository = true
	} else {
		repo := viper.GetString(core.GetFlagName(c.NS, "repository"))

		arts, _, err = client.Must().RegistryClient.ArtifactsApi.RegistriesRepositoriesArtifactsGet(
			context.Background(), regId, repo).Execute()
		if err != nil {
			return err
		}
	}

	// When listing all artifacts, add "Repository" to default cols if user hasn't specified cols
	if addRepository && cols == nil {
		cols = append(table.DefaultCols(allCols), "Repository")
	}

	return c.Out(table.Sprint(allCols, arts, cols, table.WithPrefix("items")))
}
