package provider

import (
	"context"

	cert "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func ProviderListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "provider",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "List all providers",
		LongDesc:  "List all registered ACME certificate providers across every region your account can reach. Each row shows the provider's ID, name, contact email, ACME server URL, and state. Auto-certificates reference these providers by ID.",
		Example:   "ionosctl certmanager provider list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
				certClient := cert.NewAPIClient(cfg)
				ls, _, err := certClient.ProviderApi.ProvidersGet(context.Background()).Execute()
				return ls, err
			})
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
