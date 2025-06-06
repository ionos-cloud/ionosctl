package registry

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/ionosctl/v6/services/container-registry/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{
		"RegistryId", "DisplayName", "Location", "Hostname", "VulnerabilityScanning",
		"GarbageCollectionDays", "GarbageCollectionTime", "State",
	}
)

func RegistryCmd() *core.Command {
	regCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "registry",
			Aliases: []string{"reg", "registries", "r"},
			Short:   "Registry Operations",
			Long: "Manage container registries for storage of docker images and OCI compliant artifacts. " +
				"This operation is restricted to contract owner, admin, and users with 'accessAndManageRegistries' and " +
				"Share/Edit access permissions for the data center hosting the registry.",
			TraverseChildren: true,
		},
	}

	regCmd.AddCommand(RegListCmd())
	regCmd.AddCommand(RegPostCmd())
	regCmd.AddCommand(RegGetCmd())
	regCmd.AddCommand(RegDeleteCmd())
	regCmd.AddCommand(RegUpdateCmd())
	regCmd.AddCommand(RegReplaceCmd())

	return regCmd
}

func RegsIds() []string {
	// client, _ := config.GetClient()
	svc := resources.NewRegistriesService(client.Must(), context.Background())
	regs, _, _ := svc.List("")
	return functional.Map(
		regs.GetItems(), func(reg containerregistry.RegistryResponse) string {
			return reg.GetId()
		},
	)
}
