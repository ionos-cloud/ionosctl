package registry

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "RegistryId", JSONPath: "id", Default: true},
	{Name: "DisplayName", JSONPath: "properties.name", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "Hostname", JSONPath: "properties.hostname", Default: true},
	{Name: "VulnerabilityScanning", JSONPath: "properties.features.vulnerabilityScanning.enabled", Default: true},
	{Name: "GarbageCollectionDays", JSONPath: "properties.garbageCollectionSchedule.days", Default: true},
	{Name: "GarbageCollectionTime", JSONPath: "properties.garbageCollectionSchedule.time", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func RegistryCmd() *core.Command {
	regCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "registry",
			Aliases: []string{"reg", "registries", "r"},
			Short:   "Create and manage registry instances",
			Long: `A registry is the top-level Container Registry resource: a private, OCI-compliant store for Docker images and OCI artifacts, reachable at a public HTTPS hostname (of the form <name>.cr.<location>.ionos.com).

Each registry has:
  - a name, globally unique across all IONOS customers and used as the hostname prefix (lowercase letters, digits and dashes; 3-63 chars; check availability first with 'container-registry name');
  - a location (e.g. de/txl) fixed at creation time — it cannot be changed afterwards;
  - a garbage-collection schedule (weekly days + a UTC time-of-day) that reclaims storage from untagged/deleted artifacts;
  - optional features, currently vulnerability scanning (a paid add-on, on by default).

After creation, log in with 'docker login <hostname>' using a token (see 'container-registry token'), then push/pull images. This resource is restricted to the contract owner, admins, and users with the 'accessAndManageRegistries' privilege.`,
			TraverseChildren: true,
		},
	}

	regCmd.AddColsFlag(allCols)

	regCmd.AddCommand(RegListCmd())
	regCmd.AddCommand(RegPostCmd())
	regCmd.AddCommand(RegGetCmd())
	regCmd.AddCommand(RegDeleteCmd())
	regCmd.AddCommand(RegUpdateCmd())
	regCmd.AddCommand(RegReplaceCmd())

	return regCmd
}

func RegsIds() []string {
	regs, _, _ := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
	return functional.Map(
		regs.GetItems(), func(reg containerregistry.RegistryResponse) string {
			return reg.GetId()
		},
	)
}
