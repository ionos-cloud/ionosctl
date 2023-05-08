package record

import (
	"context"
	"fmt"

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
		ShortDesc: "Create a record. Wiki: https://docs.ionos.com/dns-as-a-service/readme/api-how-tos/create-a-new-dns-record",
		Example:   "ionosctl dns record create --type A --content 1.2.3.4 --name *",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagName)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagZoneId)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagContent)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagType)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := ionoscloud.RecordProperties{}
			if fn := core.GetFlagName(c.NS, constants.FlagEnabled); viper.IsSet(fn) {
				input.Enabled = pointer.From(viper.GetBool(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Name = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagContent); viper.IsSet(fn) {
				input.Content = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagTtl); true {
				input.Ttl = pointer.From(viper.GetInt32(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagPriority); true {
				input.Priority = pointer.From(viper.GetInt32(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagType); viper.IsSet(fn) {
				input.Type = (*ionoscloud.RecordType)(pointer.From(viper.GetString(fn)))
			}

			rec, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsPost(context.Background(), id).
				RecordCreateRequest(ionoscloud.RecordCreateRequest{
					Properties: &input,
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the DNS record.  Provide a wildcard i.e. `*` to match requests for non-existent names under your DNS Zone name", core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.FlagEnabled, "", true, "When true - the record is visible for lookup")
	cmd.AddStringFlag(constants.FlagContent, "", "", fmt.Sprintf("The content (Record Data) for your chosen record type. For example, if --%s A, --%s should be an IPv4 IP.", constants.FlagType, constants.FlagContent), core.RequiredFlagOption())
	cmd.AddInt32Flag(constants.FlagTtl, "", 3600, "Time to live. The amount of time the record can be cached by a resolver or server before it needs to be refreshed from the authoritative DNS server")
	cmd.AddInt32Flag(constants.FlagPriority, "", 0, "Priority value is between 0 and 65535. Priority is mandatory for MX, SRV and URI record types and ignored for all other types.")
	cmd.AddSetFlag(constants.FlagType, "t", "AAAA",
		[]string{"A", "AAAA", "CNAME", "ALIAS", "MX", "NS", "SRV", "TXT", "CAA", "SSHFP", "TLSA", "SMIMEA", "DS", "HTTPS", "SVCB", "OPENPGPKEY", "CERT", "URI", "RP", "LOC"},
		"Type of DNS Record", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true

	return cmd
}
