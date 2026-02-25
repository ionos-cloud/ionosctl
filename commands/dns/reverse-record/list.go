package reverse_record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "list",
		Aliases:   []string{"ls", "l"},
		ShortDesc: "Retrieve all reverse records",
		Example:   "ionosctl dns rr list",
		CmdRun: func(c *core.CommandConfig) error {
			ls, err := Records(FilterRecordsByIp(c.NS))
			if err != nil {
				return fmt.Errorf("failed listing records: %w", err)
			}

			items, ok := ls.GetItemsOk()
			if !ok || items == nil {
				return fmt.Errorf("could not retrieve Record items")
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, ls, cols, table.WithPrefix("items")))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagIps, "i", "", "Optional filter for the IP address of the reverse record",
		core.WithCompletion(func() []string {
			return RecordsProperty(func(t dns.ReverseRecordRead) string {
				return t.Properties.Ip
			})
		}, constants.DNSApiRegionalURL, constants.DNSLocations),
	)

	return cmd
}
