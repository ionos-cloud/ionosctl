package record

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	dns "github.com/ionos-cloud/sdk-go-dnsaas"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

var (
	zoneId   string
	recordId string
)

func ZonesRecordsPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "update",
		Aliases:   []string{},
		ShortDesc: "Partially modify a record's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		LongDesc: fmt.Sprintf(`Partially modify a record's properties. This command uses a combination of GET and PUT to simulate a PATCH operation.
You must use either --%s and --%s, or alternatively use filters: --%s and/or --%s. Note that if choosing to use filters, the operation will fail if more than one record is found`, constants.FlagZoneId, constants.FlagRecordId, constants.FlagName, constants.FlagZoneId),
		Example: "ionosctl dns zone update --zone-id ZONE_ID --record-id RECORD_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			// Either the user directly provides the zone & record ID necessary for a findByRecordId
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagZoneId, constants.FlagRecordId)

			// Or the user provides enough filters (e.g. filter by Name or filter by ZoneId) to perform a GET /records,
			// with the mention that if more than one record is found with the given filters then we throw an error
			if f1 := core.GetFlagName(c.NS, constants.FlagRecordId); !viper.IsSet(f1) {
				err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.FlagName}, []string{constants.FlagZoneId})
				if err != nil {
					return fmt.Errorf("either provide --%s and --%s, or enough filters to narrow down to one record: --%s and/or --%s: %w", constants.FlagZoneId, constants.FlagRecordId, constants.FlagName, constants.FlagZoneId, err)
				}
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			if fn := core.GetFlagName(c.NS, constants.FlagRecordId); viper.IsSet(fn) {
				// In this case we know for sure that FlagZoneId is also set, because of the pre-run checks.
				r, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(), zoneId, recordId).Execute()
				if err != nil {
					return fmt.Errorf("failed finding record: %w", err)
				}
				return partiallyUpdateRecord(c, r)
			}

			recs, err := Records(func(r dns.ApiRecordsGetRequest) dns.ApiRecordsGetRequest {
				if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
					r = r.FilterName(viper.GetString(fn))
				}
				if fn := core.GetFlagName(c.NS, constants.FlagZoneId); viper.IsSet(fn) {
					r = r.FilterZoneId(viper.GetString(fn))
				}
				return r
			})
			if err != nil {
				return fmt.Errorf("failed listing records: %w", err)
			}

			if len(*recs.Items) > 1 {
				recsNames := functional.Fold(*recs.Items, func(acc []string, t dns.RecordResponse) []string {
					return append(acc, *t.Properties.Name)
				}, []string{})

				return fmt.Errorf("found too many records for the given filters: %+v. "+
					"The given filters (--%s and/or --%s) must narrow down to a single result",
					strings.Join(recsNames, ", "), constants.FlagName, constants.FlagZoneId)
			}

			return partiallyUpdateRecord(c, (*recs.Items)[0])

		},
		InitClient: true,
	})

	cmd.AddStringVarFlag(&zoneId, constants.FlagZoneId, "", "", "The ID (UUID) of the DNS zone", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZoneId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.ZoneIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringVarFlag(&recordId, constants.FlagRecordId, constants.FlagIdShort, "", "The ID (UUID) of the DNS record", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZoneId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RecordIds(), cobra.ShellCompDirectiveNoFileComp
	})
	return addRecordCreateFlags(cmd)
}

func partiallyUpdateRecord(c *core.CommandConfig, r dns.RecordResponse) error {
	input := r.Properties
	modifyRecordPropertiesFromFlags(c, input)

	rNew, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsPut(context.Background(), *r.Metadata.ZoneId, *r.Id).
		RecordUpdateRequest(dns.RecordUpdateRequest{Properties: input}).Execute()
	if err != nil {
		return err
	}
	return c.Printer.Print(getRecordPrint(c, rNew))
}
