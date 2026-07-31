package quota

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "quota",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Show your DNS quotas and usage",
		LongDesc:  "Show your Cloud DNS resource limits (zones, records, …) and how much of each you currently use.",
		Example:   "ionosctl dns quota get",
		CmdRun: func(c *core.CommandConfig) error {
			q, _, err := client.Must().DnsClient.QuotaApi.QuotaGet(context.Background()).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(q)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
