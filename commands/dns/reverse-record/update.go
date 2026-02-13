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
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
)

func Update() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "reverse-record",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a record",
		Example:   "ionosctl dns rr update --record OLD_RECORD_IP --name mail.example.com --ip 5.6.7.8",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagRecord); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			record, err := c.Command.Command.Flags().GetString(constants.FlagRecord)
			if err != nil {
				return err
			}
			id, err := Resolve(record)
			if err != nil {
				return fmt.Errorf("can't resolve IP to a record ID: %s", err)
			}

			r, _, err := client.Must().DnsClient.ReverseRecordsApi.ReverserecordsFindById(context.Background(), id).Execute()
			if err != nil {
				return fmt.Errorf("failed querying for reverse record ID %s: %s", id, err)
			}

			name, err := c.Command.Command.Flags().GetString(constants.FlagName)
			if err != nil {
				return err
			}
			r.Properties.Name = name

			ip, err := c.Command.Command.Flags().GetString(constants.FlagIp)
			if err != nil {
				return err
			}
			r.Properties.Ip = ip

			description, err := c.Command.Command.Flags().GetString(constants.FlagDescription)
			if err != nil {
				return err
			}
			r.Properties.Description = pointer.From(description)

			rec, _, err := client.Must().DnsClient.ReverseRecordsApi.ReverserecordsPut(context.Background(), r.Id).
				ReverseRecordEnsure(
					ionoscloud.ReverseRecordEnsure{
						Properties: r.Properties,
					}).
				Execute()
			if err != nil {
				return fmt.Errorf("failed updating record: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DnsReverseRecord, rec,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagRecord, "", "", "The record ID or IP which you want to update",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return RecordsProperty(func(read ionoscloud.ReverseRecordRead) string {
				return read.Properties.Ip
			})
		}, constants.DNSApiRegionalURL, constants.DNSLocations),
	)

	cmd.AddStringFlag(constants.FlagIp, "", "", "The new IP", core.WithCompletionE(
		func() ([]string, error) {
			ipblocks, _, err := client.Must().CloudClient.IPBlocksApi.IpblocksGet(context.Background()).Execute()
			if err != nil || ipblocks.Items == nil || len(*ipblocks.Items) == 0 {
				return nil, fmt.Errorf("failed to get IP blocks: %s", err)
			}
			var ips []string
			for _, ipblock := range *ipblocks.Items {
				if ipblock.Properties.Ips != nil {
					ips = append(ips, *ipblock.Properties.Ips...)
				}
			}
			return ips, nil
		}, "", nil),
	)
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The new record name")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "The new description of the record")

	cmd.Command.SilenceUsage = true
	return cmd
}
