package record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	dns "github.com/ionos-cloud/sdk-go-dnsaas"
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
		Example: `ionosctl dns record delete --zone-id ZONE --record-id RECORD
ionosctl dns record delete --all [--name PARTIAL_NAME] [--zone-id ZONE_ID]
ionosctl dns record delete --name PARTIAL_NAME [--zone-id ZONE_ID]`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsMutuallyExclusive(constants.ArgAll, constants.FlagRecordId)

			err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, // All with optional filters
				[]string{constants.FlagZoneId, constants.FlagRecordId},       // Known IDs
				[]string{constants.FlagName}, []string{constants.FlagZoneId}, // If none of the above, user can narrow down to a single record using filters. If more than one result, throw err
			)
			if err != nil {
				return fmt.Errorf("either provide --%s and optionally filters, or --%s and --%s, or narrow down to one record with --%s and/or --%s: %w",
					constants.ArgAll, constants.FlagZoneId, constants.FlagRecordId, constants.FlagName, constants.FlagZoneId, err)
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			r := dns.RecordResponse{}
			var err error
			if fn := core.GetFlagName(c.NS, constants.FlagRecordId); viper.IsSet(fn) {
				// In this case we know for sure that FlagZoneId is also set, because of the pre-run check
				r, _, err = client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(),
					viper.GetString(core.GetFlagName(c.NS, constants.FlagZoneId)),
					viper.GetString(core.GetFlagName(c.NS, constants.FlagRecordId)),
				).Execute()
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

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all records. Required or --%s and --%s", constants.FlagZoneId, constants.FlagRecordId))
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "If --all is set, filter --all deletion by record name")
	cmd.AddStringFlag(constants.FlagZoneId, "", "", "The zone of the target record. If --all is set, filter --all deletion by this zone id")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZoneId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.Zones(func(t dns.ZoneResponse) string {
			return *t.GetId()
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagRecordId, constants.FlagIdShort, "", fmt.Sprintf("The ID (UUID) of the DNS record. Required together with --%s or -%s", constants.FlagZoneId, constants.ArgAllShort))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRecordId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RecordIds(func(r dns.ApiRecordsGetRequest) (dns.ApiRecordsGetRequest, error) {
			if fn := core.GetFlagName(cmd.NS, constants.FlagName); viper.IsSet(fn) {
				r = r.FilterName(viper.GetString(fn))
			}
			if fn := core.GetFlagName(cmd.NS, constants.FlagZoneId); viper.IsSet(fn) {
				r = r.FilterZoneId(viper.GetString(fn))
			}
			return r, nil
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "Skip yes/no confirmation")

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	req := client.Must().DnsClient.RecordsApi.RecordsGet(c.Context)

	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		req = req.FilterName(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagZoneId); viper.IsSet(fn) {
		req = req.FilterZoneId(viper.GetString(fn))
	}

	xs, _, err := req.Execute()
	if err != nil {
		return err
	}

	err = functional.ApplyAndAggregateErrors(*xs.GetItems(), func(r dns.RecordResponse) error {
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
