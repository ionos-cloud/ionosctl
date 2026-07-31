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

Resource model:
  registry        The registry instance. Has a globally-unique name, a location, a garbage-collection schedule, and optional features (vulnerability scanning). Log in with 'docker login <hostname>'.
    token         Credentials used to authenticate 'docker push'/'docker pull' against the registry. A token has an expiry date, an enabled/disabled status, and a set of scopes.
      scope       A per-token grant: which actions (e.g. pull, push, delete) are allowed on which resource (a repository, a namespace, or the whole registry).
    repository    A named collection of images/artifacts inside a registry (the path in <hostname>/<repository>:<tag>). Created implicitly on first push.
      artifact    A single pushed image, manifest, or OCI artifact, addressed by its content digest (sha256:...).
        vulnerability  A finding from a vulnerability scan of an artifact (requires the vulnerabilityScanning feature).
  location        The set of locations (e.g. de/txl) in which a registry can be created.
  name            Checks whether a desired registry name is still available.

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
