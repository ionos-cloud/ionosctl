package record

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	ionoscloud "github.com/ionos-cloud/sdk-go-dnsaas"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

var id string

func ZonesRecordsPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a record",
		Example:   "ionosctl dns record create ",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			/* TODO: Delete/modify me for --all
						 * err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.Flag<Parent>Id}, []string{constants.ArgAll, constants.Flag<Parent>Id})
						 * if err != nil {
						 * 	return err
						 * }
			             * */

			// TODO: If no --all, mark individual flags as required

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := ionoscloud.RecordProperties{}
			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Name = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagContent); viper.IsSet(fn) {
				input.Content = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagEnabled); viper.IsSet(fn) {
				input.Enabled = pointer.From(viper.GetBool(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagEnabled); viper.IsSet(fn) {
				input.Enabled = pointer.From(viper.GetBool(fn))
			}

			rec, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsPost(context.Background(), id).
				RecordCreateRequest(ionoscloud.RecordCreateRequest{
					Properties: &ionoscloud.RecordProperties{
						Name: pointer.From("a"),
					},
				}).Execute()
			if err != nil {
				return err
			}

			return c.Printer.Print(getRecordPrint(c, rec))
		},
		InitClient: true,
	})

	cmd.AddStringVarFlag(&id, constants.FlagZoneId, constants.FlagIdShort, "", "The ID (UUID) of the DNS zone", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZoneId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.Zones(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
