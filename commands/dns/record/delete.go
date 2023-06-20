package record

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesRecordsDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a record",
		LongDesc: `To delete a specific record from a specific zone:
ionosctl dns r del --zone ZONE --record RECORD
Here, ZONE is the ID or name of the DNS zone from where you want to delete a record, and RECORD is the ID or full name of the DNS record you want to delete.

To delete all records, optionally filtering by partial name and zone:
ionosctl dns r delete --all [--record PARTIAL_NAME] [--zone ZONE]
Here, --all deletes all DNS records. You can also filter the records to delete by providing a PARTIAL_NAME that matches part of the name of the records you want to delete. Additionally, you can specify a ZONE to restrict the deletion to a specific DNS zone.

To delete a record by partial name, specifying the zone:
ionosctl dns r delete --record PARTIAL_NAME --zone ZONE
Here, PARTIAL_NAME is a part of the name of the DNS record you want to delete. If multiple records match the partial name, an error will be thrown: you will need to narrow down to a single record`,
		Example: `ionosctl dns r del --zone ZONE --record RECORD
ionosctl dns r delete --all [--record PARTIAL_NAME] [--zone ZONE]
ionosctl dns r delete --record PARTIAL_NAME --zone ZONE`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsMutuallyExclusive(constants.ArgAll, constants.FlagRecord)

			err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, // All with optional filters
				[]string{constants.FlagZone, constants.FlagRecord}, // Known resources
			)
			if err != nil {
				return fmt.Errorf("either provide --%s and optionally filters, or --%s and --%s, or narrow down to one record with --%s and/or --%s: %w",
					constants.ArgAll, constants.FlagZone, constants.FlagRecord, constants.FlagName, constants.FlagZone, err)
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			zoneId, err := zone.Resolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagZone)))
			if err != nil {
				return err
			}

			r := dns.RecordRead{}

			if fn := core.GetFlagName(c.NS, constants.FlagRecord); viper.IsSet(fn) {
				if _, ok := uuid.Parse(viper.GetString(fn)); ok != nil /* not ok (name is provided) */ {
					r, err = deleteSingleWithFilters(c)
					if err != nil {
						return fmt.Errorf("failed deleting a single record using filters: %w", err)
					}
				} else {
					r, _, err = client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(), zoneId, viper.GetString(fn)).Execute()
					if err != nil {
						return fmt.Errorf("failed finding record using Zone and Record IDs: %w", err)
					}
				}
			}

			yes := confirm.Ask(fmt.Sprintf("Are you sure you want to delete record %s (type: %s; content: %s)", *r.Metadata.Fqdn, *r.Properties.Type, *r.Properties.Content),
				viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)))
			if !yes {
				return fmt.Errorf("user cancelled deletion")
			}

			_, err = client.Must().DnsClient.RecordsApi.ZonesRecordsDelete(context.Background(),
				*r.Metadata.ZoneId,
				*r.Id,
			).Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all records. You can optionally filter the deleted records using --%s (full name / ID) and --%s (partial name)", constants.FlagZone, constants.FlagRecord))
	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", fmt.Sprintf("The full name or ID of the zone of the containing the target record. If --%s is set this is applied as a filter - limiting to records within this zone", constants.ArgAll))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.ZonesProperty(func(t dns.ZoneRead) string {
			return *t.Properties.ZoneName
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagRecord, constants.FlagRecordShort, "", fmt.Sprintf("The ID, or full name of the DNS record. Required together with --%s. Can also provide partial names, but must narrow down to a single record result if not using --%s. If using it, will however delete all records that match.", constants.FlagZone, constants.ArgAll))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRecord, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RecordsProperty(func(r dns.RecordRead) string {
			return *r.Properties.Name
		}, FilterRecordsByZoneAndRecordFlags(cmd.NS)), cobra.ShellCompDirectiveNoSpace
	})

	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "Skip yes/no confirmation")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	recordsGet := client.Must().DnsClient.RecordsApi.RecordsGet(c.Context)

	recordsGet, err := FilterRecordsByZoneAndRecordFlags(c.NS)(recordsGet)
	if err != nil {
		return err
	}

	xs, _, err := recordsGet.Execute()
	if err != nil {
		return err
	}

	err = functional.ApplyAndAggregateErrors(*xs.GetItems(), func(r dns.RecordRead) error {
		yes := confirm.Ask(fmt.Sprintf("Are you sure you want to delete record %s (content: %s)", *r.Properties.Name, *r.Properties.Content),
			viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)))

		if yes {
			_, delErr := client.Must().DnsClient.RecordsApi.ZonesRecordsDelete(c.Context, *r.Metadata.ZoneId, *r.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s (name: %s): %w", *r.Id, *r.Properties.Name, delErr)
			}
		}
		return nil
	})

	return err
}

func deleteSingleWithFilters(c *core.CommandConfig) (dns.RecordRead, error) {
	recs, err := Records(FilterRecordsByZoneAndRecordFlags(c.NS)) // full zone name and partial record name filter, if set
	if err != nil {
		return dns.RecordRead{}, fmt.Errorf("failed listing records: %w", err)
	}

	recsLen := len(*recs.Items)
	if recsLen == 0 {
		return dns.RecordRead{}, fmt.Errorf("got %d but expected 1. The given filters (--%s and/or --%s)"+
			" must narrow down to a single result", recsLen, constants.FlagRecord, constants.FlagZone)
	}

	if recsLen > 1 {
		recsNames := functional.Fold(*recs.Items, func(acc []string, t dns.RecordRead) []string {
			return append(acc, *t.Properties.Name)
		}, []string{})

		return dns.RecordRead{}, fmt.Errorf("got %d but expected 1: %+v. The given filters (--%s and/or --%s) "+
			"must narrow down to a single result. You can delete all of them by using --%s",
			recsLen, strings.Join(recsNames, ", "), constants.FlagRecord, constants.FlagZone, constants.ArgAll)
	}

	return (*recs.Items)[0], nil
}
