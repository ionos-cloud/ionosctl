package registry

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/services/container-registry/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
)

var (
	allJSONPaths = map[string]string{
		"RegistryId":            "id",
		"DisplayName":           "properties.name",
		"Location":              "properties.location",
		"Hostname":              "properties.hostname",
		"GarbageCollectionDays": "properties.garbageCollectionSchedule.days",
		"GarbageCollectionTime": "properties.garbageCollectionSchedule.time",
	}

	postCols = []string{"DisplayName", "Location"}
	allCols  = []string{"RegistryId", "DisplayName", "Location", "Hostname", "GarbageCollectionDays", "GarbageCollectionTime"}
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

	regCmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	regCmd.AddCommand(RegListCmd())
	regCmd.AddCommand(RegPostCmd())
	regCmd.AddCommand(RegGetCmd())
	regCmd.AddCommand(RegDeleteCmd())
	regCmd.AddCommand(RegUpdateCmd())
	regCmd.AddCommand(RegReplaceCmd())

	return regCmd
}

func RegsIds() []string {
	//client, _ := config.GetClient()
	svc := resources.NewRegistriesService(client.Must(), context.Background())
	regs, _, _ := svc.List("")
	return functional.Map(
		*regs.GetItems(), func(reg ionoscloud.RegistryResponse) string {
			return *reg.GetId()
		},
	)
}
