package reverse_record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Update() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "reverse-record",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Delete a record",
		Example: "ionosctl dns rr delete -af\n" +
			"ionosctl dns rr delete --record RECORD_IP\n" +
			"ionosctl dns rr delete --record RECORD_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagRecord}, []string{constants.ArgAll}); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			id, err := Resolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagRecord)))
			if err != nil {
				return fmt.Errorf("can't resolve IP to a record ID: %s", err)
			}

			r, _, err := client.Must().DnsClient.ReverseRecordsApi.ReverserecordsFindById(context.Background(), id).Execute()
			if err != nil {
				return fmt.Errorf("failed querying for reverse record ID %s: %s", id, err)
			}

			r.Properties.Name = pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
			r.Properties.Ip = pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagIp)))
			r.Properties.Description = pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagDescription)))

			_, _, err = client.Must().DnsClient.ReverseRecordsApi.ReverserecordsPut(context.Background(), *r.Id).
				ReverseRecordEnsure(
					ionoscloud.ReverseRecordEnsure{
						Properties: r.Properties,
					}).
				Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagRecord, "", "", "The record ID or IP, for identifying which record you want to update", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagRecord, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ips := RecordsProperty(func(read ionoscloud.ReverseRecordRead) string {
			return *read.Properties.Ip
		})
		return ips, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagIp, "", "", "The new IP")
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagIp, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ipblocks, _, err := client.Must().CloudClient.IPBlocksApi.IpblocksGet(context.Background()).Execute()
		if err != nil || ipblocks.Items == nil || len(*ipblocks.Items) == 0 {
			return nil, cobra.ShellCompDirectiveError
		}
		var ips []string
		for _, ipblock := range *ipblocks.Items {
			if ipblock.Properties.Ips != nil {
				ips = append(ips, *ipblock.Properties.Ips...)
			}
		}
		return ips, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagName, "", "", "The new record name")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "The new description of the record")
	// Completions
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all records if set", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	return cmd
}
