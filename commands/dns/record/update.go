package record

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
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
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a record's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		LongDesc: fmt.Sprintf(`Partially modify a record's properties. This command uses a combination of GET and PUT to simulate a PATCH operation.
You must use either --%s and --%s, or alternatively use filters: --%s and/or --%s. Note that if choosing to use filters, the operation will fail if more than one record is found`, constants.FlagZoneId, constants.FlagRecordId, constants.FlagName, constants.FlagZoneId),
		Example: "ionosctl dns zone update --zone-id ZONE_ID --record-id RECORD_ID",
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
			r, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(), zoneId, recordId).Execute()
			if err != nil {
				return fmt.Errorf("failed finding record: %w", err)
			}
			return partiallyUpdateRecordAndPrint(c, r)
		},
		InitClient: true,
	})

	cmd.AddStringVarFlag(&zoneId, constants.FlagZoneId, "", "", "The ID (UUID) of the DNS zone", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZoneId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.ZoneIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringVarFlag(&recordId, constants.FlagRecordId, constants.FlagIdShort, "", "The ID (UUID) of the DNS record", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRecordId, func(cobraCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RecordIds(func(req dns.ApiRecordsGetRequest) dns.ApiRecordsGetRequest {
			if fn := core.GetFlagName(cmd.NS, constants.FlagZoneId); viper.IsSet(fn) {
				req = req.FilterZoneId(viper.GetString(fn))
			}
			return req
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

func findRecordByListAndFilters(c *core.CommandConfig) (dns.RecordResponse, error) {
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
		return dns.RecordResponse{}, fmt.Errorf("failed listing records: %w", err)
	}

	if len(*recs.Items) > 1 {
		recsNames := functional.Fold(*recs.Items, func(acc []string, t dns.RecordResponse) []string {
			return append(acc, *t.Properties.Name)
		}, []string{})

		return dns.RecordResponse{}, fmt.Errorf("found too many records matching the given filters: %+v. "+
			"The given filters (--%s and/or --%s) must narrow down to a single result",
			strings.Join(recsNames, ", "), constants.FlagName, constants.FlagZoneId)
	}

	if len(*recs.Items) == 0 {
		return dns.RecordResponse{}, fmt.Errorf("found no records matching the given filters. "+
			"The given filters (--%s and/or --%s) must narrow down to a single result", constants.FlagName, constants.FlagZoneId)
	}

	return (*recs.Items)[0], nil
}
