package reverse_record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	dns "github.com/ionos-cloud/sdk-go-dns"
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
			ls, err := Records(FilterLimitOffset(c.NS), FilterRecordsByIp(c.NS))
			if err != nil {
				return fmt.Errorf("failed listing records: %w", err)
			}

			items, ok := ls.GetItemsOk()
			if !ok || items == nil {
				return fmt.Errorf("could not retrieve Record items")
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("items", jsonpaths.DnsReverseRecord,
				ls, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagIps, "i", "", "Optional filter for the IP address of the reverse record",
		core.WithCompletion(func() []string {
			return RecordsProperty(func(t dns.ReverseRecordRead) string {
				return *t.Properties.Ip
			})
		}, constants.DNSApiRegionalURL, constants.DNSLocations),
	)
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "The first element (of the total list of elements) to include in the response. Use together with limit for pagination")
	cmd.AddInt32Flag(constants.FlagMaxResults, "", 0, constants.DescMaxResults)

	return cmd
}
