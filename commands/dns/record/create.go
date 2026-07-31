package record

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"

	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

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
		ShortDesc: "Create a DNS record",
		LongDesc: `Create a DNS record inside a zone.

Three things define a record: --type, --name and --content. --name is the host under the zone ('www' for www.example.com, '@' or the zone name for the apex, '*' for a wildcard). --content is the record's data and its meaning depends on --type:

  A       IPv4 address            e.g. 1.2.3.4
  AAAA    IPv6 address            e.g. 2001:db8::1
  CNAME   target hostname         e.g. www.example.com
  ALIAS   target hostname (apex)  e.g. example.com
  MX      mail server hostname    e.g. mail.example.com   (set --priority)
  NS      name server hostname    e.g. ns1.example.com
  TXT     free text               e.g. "v=spf1 -all"
  SRV     "weight port target"    e.g. "5 5060 sip.example.com"  (set --priority)
  CAA     "flags tag value"       e.g. "0 issue letsencrypt.org"

--priority is required for MX, SRV and URI and ignored otherwise. --ttl sets the cache lifetime in seconds (60-604800, default 3600). Records are --enabled by default.`,
		Example: `ionosctl dns record create --zone example.com --type A --name www --content 1.2.3.4
ionosctl dns record create --zone example.com --type MX --name @ --content mail.example.com --priority 10 --ttl 300`,
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
					Properties: input,
				}).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(rec)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", "The ID or name of the DNS zone", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ZonesProperty(func(t dns.ZoneRead) string {
			return t.Properties.ZoneName
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return addRecordCreateFlags(cmd)
}

func addRecordCreateFlags(cmd *core.Command) *core.Command {
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Host under the zone this record answers for, e.g. 'www', '@' for the apex, or '*' for a wildcard matching non-existent names. Some shells need '*' escaped as '\\*'", core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.FlagEnabled, "", true, "Whether the record answers lookups. true = live; false = kept but not served (default true)")
	cmd.AddStringFlag(constants.FlagContent, "", "", fmt.Sprintf("Record data, interpreted per --%s: an A record takes an IPv4 (1.2.3.4), AAAA an IPv6, CNAME/MX/NS a hostname, TXT free text. See this command's --help for the full per-type table", constants.FlagType), core.RequiredFlagOption())
	cmd.AddInt32Flag(constants.FlagTtl, "", 3600, "How long (seconds) resolvers may cache this record before re-querying; 60-604800 (default 3600 = 1h)")
	cmd.AddInt32Flag(constants.FlagPriority, "", 0, "Preference value 0-65535, lower wins. Required for MX, SRV and URI records; ignored for all other types")
	cmd.AddSetFlag(constants.FlagType, "", "AAAA",
		[]string{"A", "AAAA", "CNAME", "ALIAS", "MX", "NS", "SRV", "TXT", "CAA", "SSHFP", "TLSA", "SMIMEA", "DS", "HTTPS", "SVCB", "OPENPGPKEY", "CERT", "URI", "RP", "LOC"},
		"Record type; decides how --content is interpreted (A=IPv4, AAAA=IPv6, CNAME/MX/NS=hostname, TXT=text, …)", core.RequiredFlagOption())

	return cmd
}

func modifyRecordPropertiesFromFlags(c *core.CommandConfig, input *dns.Record) {
	if fn := core.GetFlagName(c.NS, constants.FlagEnabled); viper.IsSet(fn) {
		input.Enabled = pointer.From(viper.GetBool(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		input.Name = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagContent); viper.IsSet(fn) {
		input.Content = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagTtl); true {
		input.Ttl = pointer.From(viper.GetInt32(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagPriority); true {
		input.Priority = pointer.From(viper.GetInt32(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagType); viper.IsSet(fn) {
		input.Type = dns.RecordType(viper.GetString(fn))
	}
}
