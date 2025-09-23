package container_registry

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/artifacts"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/location"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/name"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/token"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/vulnerabilities"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func ContainerRegistryCmd() *core.Command {
	contregCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "container-registry",
			Short:   "Registry Operations",
			Aliases: []string{"cr", "contreg", "cont-reg"},
			Long: "The sub-commands of `ionosctl container-registry` allows IONOS clients to manage docker and OCI compliant " +
				"registries for use by their managed Kubernetes clusters. Use a Container Registry to ensure you have a " +
				"privately accessed registry to efficiently support image pulls.",
			TraverseChildren: true,
		},
	}

	contregCmd.AddCommand(registry.RegistryCmd())
	contregCmd.AddCommand(token.TokenCmd())
	contregCmd.AddCommand(location.RegLocationsListCmd())
	contregCmd.AddCommand(name.RegNamesCmd())
	contregCmd.AddCommand(repository.RegRepoDeleteCmd())
	contregCmd.AddCommand(artifacts.ArtifactsCmd())
	contregCmd.AddCommand(vulnerabilities.VulnerabilitiesCmd())

	return core.WithConfigOverride(contregCmd, []string{fileconfiguration.ContainerRegistry}, constants.DefaultApiURL+"/containerregistries")
}
