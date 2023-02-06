package container_registry

import (
	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"
)

func ContainerRegistryCmd() *core.Command {
	contregCmd := &core.Command{
		Command: &cobra.Command{
			Use:   "container-registry",
			Short: "Registry Operations",
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

	return contregCmd
}
