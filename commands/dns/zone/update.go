package zone

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func ZonesPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Update a primary DNS zone",
		LongDesc: `Partially update a primary DNS zone. Only the flags you pass change; the rest are preserved (a GET+PUT that simulates PATCH). Identify the zone by name or ID with --zone.

Common use: --enabled=false to take a zone out of service without deleting its records, or --enabled=true to bring it back.`,
		Example: `ionosctl dns zone update --zone example.com --description "moved to prod"
ionosctl dns zone update --zone example.com --enabled=false`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			id, err := utils.ZoneResolve(viper.GetString(core.GetFlagName(c.NS, constants.FlagZone)))
			if err != nil {
				return err
			}

			z, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), id).Execute()
			if err != nil {
				return err
			}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				z.Properties.ZoneName = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
				z.Properties.Description = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagEnabled); viper.IsSet(fn) {
				z.Properties.Enabled = pointer.From(viper.GetBool(fn))
			}

			zNew, _, err := client.Must().DnsClient.ZonesApi.ZonesPut(context.Background(), id).
				ZoneEnsure(dns.ZoneEnsure{Properties: z.Properties}).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(zNew)
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
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "New domain name for the zone, e.g. example.com")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "New free-text note; not served in DNS")
	cmd.AddBoolFlag(constants.FlagEnabled, "", true, "Whether the zone is served. true = IONOS answers lookups; false = kept but not resolved")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
