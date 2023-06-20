package record

import (
	"context"
	"fmt"
	"strings"

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
		Example: `ionosctl dns r delete --zone ZONE --record-id RECORD
ionosctl dns r delete --all [--name PARTIAL_NAME] [--zone ZONE]
ionosctl dns r delete --name PARTIAL_NAME [--zone ZONE]`,
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

			recordId, err := Resolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagRecord)))
			if err != nil {
				return err
			}

			r := dns.RecordRead{}
			if fn := core.GetFlagName(c.NS, constants.FlagRecord); viper.IsSet(fn) {
				// In this case we know for sure that FlagZone is also set, because of the pre-run check
				r, _, err = client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(), zoneId, recordId).Execute()
				if err != nil {
					return fmt.Errorf("failed finding record using Zone and Record IDs: %w", err)
				}
			} else {
				r, err = findRecordByListAndFilters(c)
				if err != nil {
					return fmt.Errorf("failed attempt to narrow down record using filters: %w", err)
				}
			}

			yes := confirm.Ask(fmt.Sprintf("Are you sure you want to delete record %s (type: %s; content: %s)", *r.Properties.Name, *r.Properties.Type, *r.Properties.Content),
				viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)))
			if !yes {
				return nil
			}

			_, err = client.Must().DnsClient.RecordsApi.ZonesRecordsDelete(context.Background(),
				*r.Metadata.ZoneId,
				*r.Id,
			).Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all records. Required or --%s and --%s", constants.FlagZone, constants.FlagRecord))
	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", "The zone of the target record. If --all is set, filter --all deletion by limiting to records within this zone")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.Zones(func(t dns.ZoneRead) string {
			return *t.Properties.ZoneName
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagRecord, constants.FlagRecordShort, "", fmt.Sprintf("The ID or name of the DNS record. Required together with --%s", constants.FlagZone))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRecord, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RecordIds(func(r dns.ApiRecordsGetRequest) (dns.ApiRecordsGetRequest, error) {
			if fn := core.GetFlagName(cmd.NS, constants.FlagName); viper.IsSet(fn) {
				r = r.FilterName(viper.GetString(fn))
			}
			if fn := core.GetFlagName(cmd.NS, constants.FlagZone); viper.IsSet(fn) {
				zoneId, err := zone.Resolve(viper.GetString(fn))
				if err != nil {
					return r, err
				}
				r = r.FilterZoneId(zoneId)
			}
			return r, nil
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "Skip yes/no confirmation")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	req := client.Must().DnsClient.RecordsApi.RecordsGet(c.Context)

	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		req = req.FilterName(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagZone); viper.IsSet(fn) {
		zoneId, err := zone.Resolve(viper.GetString(fn))
		if err != nil {
			return err
		}
		req = req.FilterZoneId(zoneId)
	}

	xs, _, err := req.Execute()
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

func findRecordByListAndFilters(c *core.CommandConfig) (dns.RecordRead, error) {
	recs, err := Records(func(r dns.ApiRecordsGetRequest) (dns.ApiRecordsGetRequest, error) {
		if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
			r = r.FilterName(viper.GetString(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagZone); viper.IsSet(fn) {
			zoneId, err := zone.Resolve(viper.GetString(fn))
			if err != nil {
				return r, err
			}
			r = r.FilterZoneId(zoneId)
		}
		return r, nil
	})
	if err != nil {
		return dns.RecordRead{}, fmt.Errorf("failed listing records: %w", err)
	}

	if len(*recs.Items) > 1 {
		recsNames := functional.Fold(*recs.Items, func(acc []string, t dns.RecordRead) []string {
			return append(acc, *t.Properties.Name)
		}, []string{})

		return dns.RecordRead{}, fmt.Errorf("found too many records matching the given filters: %+v. "+
			"The given filters (--%s and/or --%s) must narrow down to a single result",
			strings.Join(recsNames, ", "), constants.FlagName, constants.FlagZone)
	}

	if len(*recs.Items) == 0 {
		return dns.RecordRead{}, fmt.Errorf("found no records matching the given filters. "+
			"The given filters (--%s and/or --%s) must narrow down to a single result", constants.FlagName, constants.FlagZone)
	}

	return (*recs.Items)[0], nil
}
