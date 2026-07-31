package artifacts

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ArtifactsGetCmd() *core.Command {
	c := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "artifacts",
			Verb:       "get",
			ShortDesc:  "Get a single artifact by digest",
			LongDesc:   "Get one artifact from a repository, identified by its content digest (sha256:...). Returns its media type, usage stats and vulnerability counts. List a repository's artifacts first with 'container-registry artifacts list' to find the digest.",
			Example:    "ionosctl container-registry artifacts get --registry-id REGISTRY_ID --repository my-app --artifact-id sha256:DIGEST",
			PreCmdRun:  PreCmdGet,
			CmdRun:     CmdGet,
			InitClient: true,
		},
	)

	c.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "The unique ID of the registry the artifact belongs to")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag("repository", "", "", "Name of the repository that holds the artifact")
	_ = c.Command.RegisterFlagCompletionFunc(
		"repository", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return repository.RepositoryNames(viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))),
				cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagArtifactId, "", "", "Content digest of the artifact, e.g. sha256:12ab... (as shown in the Id column of 'artifacts list')")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagArtifactId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ArtifactsIds(
					viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId)),
					viper.GetString(core.GetFlagName(c.NS, "repository")),
				),
				cobra.ShellCompDirectiveNoFileComp
		},
	)

	return c
}

func PreCmdGet(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagRegistryId, "repository", constants.FlagArtifactId)
}

func CmdGet(c *core.CommandConfig) error {
	regId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))
	repo := viper.GetString(core.GetFlagName(c.NS, "repository"))
	artId := viper.GetString(core.GetFlagName(c.NS, constants.FlagArtifactId))

	arts, _, err := client.Must().RegistryClient.ArtifactsApi.RegistriesRepositoriesArtifactsFindByDigest(
		c.Context, regId, repo, artId,
	).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(arts)
}
