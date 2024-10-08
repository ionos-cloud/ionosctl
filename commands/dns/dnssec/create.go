package dnssec

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	const (
		FlagAlgorithm       = "algorithm"
		FlagKskBits         = "ksk-bits"
		FlagZskBits         = "zsk-bits"
		FlagNsecMode        = "nsec-mode"
		FlagNsec3Iterations = "nsec3-iterations"
		FlagNsec3SaltBits   = "nsec3-salt-bits"
		FlagValidity        = "validity"
	)

	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "dnssec",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Enable DNSSEC keys and create associated DNSKEY records for your DNS zone",
		Example:   "ionosctl dns keys create --zone ZONE",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			zoneId, err := utils.ZoneResolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagZone)))
			if err != nil {
				return err
			}

			key, _, err := client.Must().DnsClient.DNSSECApi.ZonesKeysPost(context.Background(), zoneId).
				DnssecKeyCreate(
					dns.DnssecKeyCreate{
						Properties: &dns.DnssecKeyParameters{
							Validity: pointer.From(viper.GetInt32(core.GetFlagName(c.NS, FlagValidity))),
							KeyParameters: &dns.KeyParameters{
								Algorithm: pointer.From(dns.Algorithm(viper.GetString(core.GetFlagName(c.NS, FlagAlgorithm)))),
								KskBits:   pointer.From(dns.KskBits(viper.GetInt32(core.GetFlagName(c.NS, FlagKskBits)))),
								ZskBits:   pointer.From(dns.ZskBits(viper.GetInt32(core.GetFlagName(c.NS, FlagZskBits)))),
							},
							NsecParameters: &dns.NsecParameters{
								NsecMode:        pointer.From(dns.NsecMode(viper.GetString(core.GetFlagName(c.NS, FlagNsecMode)))),
								Nsec3Iterations: pointer.From(viper.GetInt32(core.GetFlagName(c.NS, FlagNsec3Iterations))),
								Nsec3SaltBits:   pointer.From(viper.GetInt32(core.GetFlagName(c.NS, FlagNsec3SaltBits))),
							},
						},
					}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DnsSecKey, key,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone)
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ZonesProperty(func(t dns.ZoneRead) string {
			return *t.Properties.ZoneName
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(FlagAlgorithm, "", "RSASHA256", "Algorithm used to generate signing keys (both Key Signing Keys and Zone Signing Keys)")
	cmd.AddIntFlag(FlagKskBits, "", 1024, "Key signing key length in bits. kskBits >= zskBits: [1024/2048/4096]")
	_ = cmd.Command.RegisterFlagCompletionFunc(FlagKskBits, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1024", "2048", "4096"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(FlagZskBits, "", 1024, "Zone signing key length in bits. zskBits <= kskBits: [1024/2048/4096]")
	_ = cmd.Command.RegisterFlagCompletionFunc(FlagZskBits, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1024", "2048", "4096"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(FlagNsecMode, "", "NSEC", []string{"NSEC", "NSEC3"}, "NSEC mode.")
	cmd.AddIntFlag(FlagNsec3Iterations, "", 0, "Number of iterations for NSEC3. [0..50]")
	cmd.AddIntFlag(FlagNsec3SaltBits, "", 64, "Salt length in bits for NSEC3. [64..128], multiples of 8")
	_ = cmd.Command.RegisterFlagCompletionFunc(FlagNsec3SaltBits, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"64", "72", "80", "88", "96", "104", "112", "120", "128"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(FlagValidity, "", 90, "Signature validity in days [90..365]")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
