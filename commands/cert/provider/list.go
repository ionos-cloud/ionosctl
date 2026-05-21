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
		ShortDesc: "Retrieve Provider list",
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
