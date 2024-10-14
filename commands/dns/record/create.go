package record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"

	dns "github.com/ionos-cloud/sdk-go-dns"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func ZonesRecordsPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a record. Wiki: https://docs.ionos.com/cloud/network-services/cloud-dns/api-how-tos/create-dns-record",
		Example:   "ionosctl dns r create --zone foo-bar.com --type A --content 1.2.3.4 --name \\*",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagZone, constants.FlagContent, constants.FlagType); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := dns.Record{}
			modifyRecordPropertiesFromFlags(c, &input)

			zoneId, err := utils.ZoneResolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagZone)))
			if err != nil {
				return err
			}

			rec, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsPut(context.Background(), zoneId, uuidgen.Must()).
				RecordEnsure(dns.RecordEnsure{
					Properties: &input,
				}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			// if err != nil {
			//	return err
			// }

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DnsRecord, rec,
				tabheaders.GetHeadersAllDefault(defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", "The ID or name of the DNS zone", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ZonesProperty(func(t dns.ZoneRead) string {
			return *t.Properties.ZoneName
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return addRecordCreateFlags(cmd)
}

func addRecordCreateFlags(cmd *core.Command) *core.Command {
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the DNS record.  Provide a wildcard i.e. `\\*` to match requests for non-existent names under your DNS Zone name. Note that some terminals require '*' to be escaped, e.g. '\\*'", core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.FlagEnabled, "", true, "When true - the record is visible for lookup")
	cmd.AddStringFlag(constants.FlagContent, "", "", fmt.Sprintf("The content (Record Data) for your chosen record type. For example, if --%s A, --%s should be an IPv4 IP.", constants.FlagType, constants.FlagContent), core.RequiredFlagOption())
	cmd.AddInt32Flag(constants.FlagTtl, "", 3600, "Time to live. The amount of time the record can be cached by a resolver or server before it needs to be refreshed from the authoritative DNS server")
	cmd.AddInt32Flag(constants.FlagPriority, "", 0, "Priority value is between 0 and 65535. Priority is mandatory for MX, SRV and URI record types and ignored for all other types.")
	cmd.AddSetFlag(constants.FlagType, "t", "AAAA",
		[]string{"A", "AAAA", "CNAME", "ALIAS", "MX", "NS", "SRV", "TXT", "CAA", "SSHFP", "TLSA", "SMIMEA", "DS", "HTTPS", "SVCB", "OPENPGPKEY", "CERT", "URI", "RP", "LOC"},
		"Type of DNS Record", core.RequiredFlagOption())

	return cmd
}

func modifyRecordPropertiesFromFlags(c *core.CommandConfig, input *dns.Record) {
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
		input.Type = pointer.From(viper.GetString(fn))
	}
}
