package reverse_record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/spf13/viper"
)

func Delete() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "reverse-record",
		Verb:      "delete",
		Aliases:   []string{"d", "del", "rm"},
		ShortDesc: "Delete a record",
		Example: "ionosctl dns rr delete -af\n" +
			"ionosctl dns rr delete --record RECORD_IP\n" +
			"ionosctl dns rr delete --record RECORD_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagRecord}, []string{constants.FlagAll}); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.FlagAll)); all {
				return deleteAll(c)
			}

			return deleteSingle(c, viper.GetString(core.GetFlagName(c.NS, constants.FlagRecord)))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagRecord, "", "", "The record ID or IP which you want to delete",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return RecordsProperty(func(read ionoscloud.ReverseRecordRead) string {
				return read.Properties.Ip
			})
		}, constants.DNSApiRegionalURL, constants.DNSLocations),
	)

	cmd.AddBoolFlag(constants.FlagAll, constants.FlagAllShort, false, "Delete all records if set", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	records, err := Records()
	if err != nil {
		return fmt.Errorf("failed getting all records: %w", err)
	}

	return functional.ApplyAndAggregateErrors(records.Items, func(r ionoscloud.ReverseRecordRead) error {
		return deleteSingle(c, r.Id)
	})
}

func deleteSingle(c *core.CommandConfig, ipOrIdOfRecord string) error {
	id, err := Resolve(ipOrIdOfRecord)
	if err != nil {
		return fmt.Errorf("can't resolve IP %s to a record ID: %s", ipOrIdOfRecord, err)
	}

	r, _, err := client.Must().DnsClient.ReverseRecordsApi.ReverserecordsFindById(context.Background(), id).Execute()
	if err != nil {
		return fmt.Errorf("failed querying for reverse record ID %s: %s", id, err)
	}
	yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf(
		"Are you sure you want to delete record %s (IP: '%s'; description: '%s'; ID: '%s')",
		r.Properties.Name, r.Properties.Ip, *r.Properties.Description, r.Id),
		viper.GetBool(constants.FlagForce))
	if !yes {
		return fmt.Errorf("user cancelled deletion")
	}

	_, _, err = client.Must().DnsClient.ReverseRecordsApi.ReverserecordsDelete(context.Background(),
		r.Id,
	).Execute()

	return err
}
