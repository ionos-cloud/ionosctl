package reverse_record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/spf13/viper"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "reverse-record",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Find a record by IP or ID",
		Example: "ionosctl dns rr get --record RECORD_IP\n" +
			"ionosctl dns rr get --record RECORD_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagRecord}, []string{constants.ArgAll}); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			id, err := Resolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagRecord)))
			if err != nil {
				return fmt.Errorf("can't resolve IP to a record ID: %s", err)
			}

			rec, _, err := client.Must().DnsClient.ReverseRecordsApi.ReverserecordsFindById(context.Background(), id).Execute()
			if err != nil {
				return fmt.Errorf("failed querying for reverse record ID %s: %s", id, err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DnsReverseRecord, rec,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagRecord, "", "", "The record ID or IP which you want to get",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return RecordsProperty(func(read ionoscloud.ReverseRecordRead) string {
				return *read.Properties.Ip
			})
		}, constants.DNSApiRegionalURL, constants.DNSLocations),
	)

	cmd.Command.SilenceUsage = true
	return cmd
}
