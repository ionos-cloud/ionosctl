package container_registry

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/location"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/name"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/token"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
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

	contregCmd.Command.PersistentFlags().Bool(
		constants.ArgNoHeaders, false, "When using text output, don't print headers",
	)

	contregCmd.AddCommand(registry.RegistryCmd())
	contregCmd.AddCommand(token.TokenCmd())
	contregCmd.AddCommand(location.RegLocationsListCmd())
	contregCmd.AddCommand(name.RegNamesCmd())
	contregCmd.AddCommand(repository.RegRepoDeleteCmd())

	return contregCmd
}
