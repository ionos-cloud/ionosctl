package provider

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
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
			req := client.Must().CertManagerClient.ProviderApi.ProvidersGet(context.Background())

			ls, _, err := req.Execute()
			if err != nil {
				return fmt.Errorf("failed listing the Providers: %w", err)
			}

			return c.Printer(allCols).Prefix("items").Print(ls)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
