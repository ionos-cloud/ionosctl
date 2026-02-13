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
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
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
			zone, err := c.Command.Command.Flags().GetString(constants.FlagZone)
			if err != nil {
				return err
			}
			zoneId, err := utils.ZoneResolve(zone)
			if err != nil {
				return err
			}

			validity, err := c.Command.Command.Flags().GetInt32(FlagValidity)
			if err != nil {
				return err
			}

			algorithm, err := c.Command.Command.Flags().GetString(FlagAlgorithm)
			if err != nil {
				return err
			}

			kskBits, err := c.Command.Command.Flags().GetInt32(FlagKskBits)
			if err != nil {
				return err
			}

			zskBits, err := c.Command.Command.Flags().GetInt32(FlagZskBits)
			if err != nil {
				return err
			}

			nsecMode, err := c.Command.Command.Flags().GetString(FlagNsecMode)
			if err != nil {
				return err
			}

			nsec3Iterations, err := c.Command.Command.Flags().GetInt32(FlagNsec3Iterations)
			if err != nil {
				return err
			}

			nsec3SaltBits, err := c.Command.Command.Flags().GetInt32(FlagNsec3SaltBits)
			if err != nil {
				return err
			}

			key, _, err := client.Must().DnsClient.DNSSECApi.ZonesKeysPost(context.Background(), zoneId).
				DnssecKeyCreate(
					dns.DnssecKeyCreate{
						Properties: dns.DnssecKeyParameters{
							Validity: validity,
							KeyParameters: dns.KeyParameters{
								Algorithm: dns.Algorithm(algorithm),
								KskBits:   dns.KskBits(kskBits),
								ZskBits:   dns.ZskBits(zskBits),
							},
							NsecParameters: dns.NsecParameters{
								NsecMode:        dns.NsecMode(nsecMode),
								Nsec3Iterations: nsec3Iterations,
								Nsec3SaltBits:   nsec3SaltBits,
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

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ZonesProperty(func(t dns.ZoneRead) string {
				return t.Properties.ZoneName
			})
		}, constants.DNSApiRegionalURL, constants.DNSLocations),
	)
	cmd.AddStringFlag(FlagAlgorithm, "", "RSASHA256", "Algorithm used to generate signing keys (both Key Signing Keys and Zone Signing Keys)")
	cmd.AddIntFlag(FlagKskBits, "", 1024, "Key signing key length in bits. kskBits >= zskBits: [1024/2048/4096]",
		core.WithCompletion(
			func() []string {
				return []string{"1024", "2048", "4096"}
			}, constants.DNSApiRegionalURL, constants.DNSLocations,
		),
	)
	cmd.AddIntFlag(FlagZskBits, "", 1024, "Zone signing key length in bits. zskBits <= kskBits: [1024/2048/4096]",
		core.WithCompletion(
			func() []string {
				return []string{"1024", "2048", "4096"}
			}, constants.DNSApiRegionalURL, constants.DNSLocations,
		),
	)
	cmd.AddSetFlag(FlagNsecMode, "", "NSEC", []string{"NSEC", "NSEC3"}, "NSEC mode.")
	cmd.AddIntFlag(FlagNsec3Iterations, "", 0, "Number of iterations for NSEC3. [0..50]")
	cmd.AddIntFlag(FlagNsec3SaltBits, "", 64, "Salt length in bits for NSEC3. [64..128], multiples of 8",
		core.WithCompletion(
			func() []string {
				return []string{"64", "72", "80", "88", "96", "104", "112", "120", "128"}
			}, constants.DNSApiRegionalURL, constants.DNSLocations,
		),
	)
	cmd.AddIntFlag(FlagValidity, "", 90, "Signature validity in days [90..365]")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
