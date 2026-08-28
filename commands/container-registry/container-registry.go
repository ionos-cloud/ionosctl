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
			Short:   "Manage private Docker/OCI container registries",
			Aliases: []string{"cr", "contreg", "cont-reg"},
			Long: `The IONOS Container Registry (CR) is a private, OCI-compliant registry for storing and distributing Docker images and OCI artifacts (e.g. Helm charts). It is a regional service: each registry lives in one location and is reached over a public HTTPS endpoint (hostname), so it can be used from Managed Kubernetes clusters, CI pipelines, or any Docker/OCI client.

A registry holds repositories (created implicitly on first push), each addressing its images and artifacts by content digest. Access tokens authenticate 'docker login', each scoped to a set of actions (pull, push, delete) on a repository, a namespace, or the whole registry. Optional vulnerability scanning reports findings per artifact.

Access is restricted to the contract owner, admins, and users with the 'accessAndManageRegistries' privilege.`,
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
