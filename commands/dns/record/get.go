package record

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"
	dns "github.com/ionos-cloud/sdk-go-dnsaas"
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/viper"
)

func ZonesRecordsFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a record",
		Example:   "ionosctl dns record get --zoneId ZONE_ID --recordId RECORD_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagZoneId)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagRecordId)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			r, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagZoneId)),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagRecordId)),
			).Execute()

			if err != nil {
				return err
			}
			return c.Printer.Print(getRecordPrint(c, r))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZoneId, "", "", "The ID (UUID) of the DNS zone of which record belongs to")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZoneId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.ZoneIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagRecordId, "", "", "The ID (UUID) of the DNS record")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRecordId, func(cobraCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RecordIds(func(req dns.ApiRecordsGetRequest) dns.ApiRecordsGetRequest {
			if fn := core.GetFlagName(cmd.NS, constants.FlagZoneId); viper.IsSet(fn) {
				req = req.FilterZoneId(viper.GetString(fn))
			}
			return req
		}), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
