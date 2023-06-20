package record

import (
	"context"

	"github.com/google/uuid"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"

	dns "github.com/ionos-cloud/sdk-go-dns"
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
		Example:   "ionosctl dns r get --zone ZONE --record RECORD",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagZone)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagRecord)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			zoneId, err := zone.Resolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagZone)))
			if err != nil {
				return err
			}
			recordId, err := Resolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagRecord)))
			if err != nil {
				return err
			}
			r, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(),
				zoneId, recordId,
			).Execute()

			if err != nil {
				return err
			}
			return c.Printer.Print(getRecordPrint(c, r))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone)
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.Zones(func(t dns.ZoneRead) string {
			return *t.Properties.ZoneName
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagRecord, "", "", "The ID or name of the DNS record")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRecord, func(cobraCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RecordIds(func(req dns.ApiRecordsGetRequest) (dns.ApiRecordsGetRequest, error) {
			if fn := core.GetFlagName(cmd.NS, constants.FlagZone); viper.IsSet(fn) {
				zoneId, err := zone.Resolve(viper.GetString(fn))
				if err != nil {
					return req, err
				}
				req = req.FilterZoneId(zoneId)
			}

			if fn := core.GetFlagName(cmd.NS, constants.FlagRecord); viper.IsSet(fn) {
				record := viper.GetString(fn)
				if _, ok := uuid.Parse(record); ok == nil /* not ok (name is provided) */ {
					req = req.FilterName(record)
				}
			}
			return req, nil
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
