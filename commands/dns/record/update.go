package record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

var (
	recordId string
)

func ZonesRecordsPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a record's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		Example:   "ionosctl dns zone update --zone ZONE_ID --record-id RECORD_ID",
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
			r, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(), zoneId, recordId).Execute()
			if err != nil {
				return fmt.Errorf("failed finding record: %w", err)
			}
			return partiallyUpdateRecordAndPrint(c, r)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.Zones(func(t dns.ZoneResponse) string {
			return *t.Properties.ZoneName
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringVarFlag(&recordId, constants.FlagRecord, constants.FlagIdShort, "", "The ID (UUID) of the DNS record", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRecord, func(cobraCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RecordIds(func(req dns.ApiRecordsGetRequest) (dns.ApiRecordsGetRequest, error) {
			if fn := core.GetFlagName(cmd.NS, constants.FlagZone); viper.IsSet(fn) {
				zoneId, err := zone.Resolve(viper.GetString(fn))
				if err != nil {
					return req, err
				}
				req = req.FilterZoneId(zoneId)
			}
			return req, nil
		}), cobra.ShellCompDirectiveNoFileComp
	})
	return addRecordCreateFlags(cmd)
}

func partiallyUpdateRecordAndPrint(c *core.CommandConfig, r dns.RecordResponse) error {
	input := r.Properties
	modifyRecordPropertiesFromFlags(c, input)

	rNew, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsPut(context.Background(), *r.Metadata.ZoneId, *r.Id).
		RecordUpdateRequest(dns.RecordUpdateRequest{Properties: input}).Execute()
	if err != nil {
		return err
	}
	return c.Printer.Print(getRecordPrint(c, rNew))
}
