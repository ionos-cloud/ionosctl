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
		Example:   `ionosctl dns record delete (--zone-id ZONE --record-id RECORD | --all [--name PARTIAL_NAME] [--zone-id ZONE_ID])`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagZoneId, constants.FlagRecordId})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			z, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagZoneId)),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagRecordId)),
			).Execute()
			if err != nil {
				return fmt.Errorf("failed getting record by id %s", id)
			}
			yes := confirm.Ask(fmt.Sprintf("Are you sure you want to delete record %s (content: %s)", *z.Properties.Name, *z.Properties.Content),
				viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)))
			if !yes {
				return nil
			}

			_, err = client.Must().DnsClient.RecordsApi.ZonesRecordsDelete(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagZoneId)),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagRecordId)),
			).Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all records. Required or --%s and --%s", constants.FlagZoneId, constants.FlagRecordId))
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "If --all is set, filter --all deletion by record name")
	cmd.AddStringFlag(constants.FlagZoneId, "", "", "The zone of the target record. If --all is set, filter --all deletion by this zone id")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZoneId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return zone.ZoneIds(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagRecordId, constants.FlagIdShort, "", fmt.Sprintf("The ID (UUID) of the DNS record. Required together with --%s or -%s", constants.FlagZoneId, constants.ArgAllShort))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRecordId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RecordIds(func(r dns.ApiRecordsGetRequest) dns.ApiRecordsGetRequest {
			if fn := core.GetFlagName(cmd.NS, constants.FlagName); viper.IsSet(fn) {
				r = r.FilterName(viper.GetString(fn))
			}
			if fn := core.GetFlagName(cmd.NS, constants.FlagZoneId); viper.IsSet(fn) {
				r = r.FilterZoneId(viper.GetString(fn))
			}
			return r
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
