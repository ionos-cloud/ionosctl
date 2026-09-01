package vulnerabilities

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/artifacts"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VulnerabilitiesListCmd() *core.Command {
	c := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "vulnerabilities",
			Verb:       "list",
			ShortDesc:  "List vulnerabilities found in an artifact",
			LongDesc:   "List all vulnerability findings for a single artifact, identified by registry, repository and artifact digest. Requires the registry's vulnerabilityScanning feature to be enabled; if it is off, no findings are returned. Each row shows the CVE, score, severity and whether a fix exists.",
			Example:    "ionosctl container-registry vulnerabilities list --registry-id REGISTRY_ID --repository my-app --artifact-id sha256:DIGEST",
			PreCmdRun:  PreCmdList,
			CmdRun:     CmdList,
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

	c.AddStringFlag(constants.FlagArtifactId, "", "", "Content digest of the artifact to scan for findings, e.g. sha256:12ab...")
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

	return c
}

func PreCmdList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagRegistryId, "repository", constants.FlagArtifactId,
	); err != nil {
		return err
	}

	return nil
}

func CmdList(c *core.CommandConfig) error {
	registryId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))
	repository := viper.GetString(core.GetFlagName(c.NS, "repository"))
	artifactId := viper.GetString(core.GetFlagName(c.NS, constants.FlagArtifactId))

	vulnerabilities, _, err := client.Must().RegistryClient.ArtifactsApi.
		RegistriesRepositoriesArtifactsVulnerabilitiesGet(
			context.Background(), registryId, repository, artifactId).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(vulnerabilities)
}
