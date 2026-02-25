package quota

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "quota",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve your quotas",
		Example:   "ionosctl dns quota get",
		CmdRun: func(c *core.CommandConfig) error {
			q, _, err := client.Must().DnsClient.QuotaApi.QuotaGet(context.Background()).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return table.Fprint(c.Command.Command.OutOrStdout(), allCols, q, cols)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
