package zone

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func ZonesPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Create a primary DNS zone",
		LongDesc: `Create a primary DNS zone for a domain you want IONOS to answer for.

--name is the domain itself (e.g. example.com), NOT a friendly label. After creating the zone, delegate the domain to the IONOS name servers at your registrar and add entries with 'dns record create'. A zone starts --enabled; pass --enabled=false to create it dormant.`,
		Example: `ionosctl dns zone create --name example.com
ionosctl dns zone create --name example.com --description "prod apex" --enabled=false`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := dns.Zone{}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.ZoneName = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
				input.Description = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagEnabled); viper.IsSet(fn) {
				input.Enabled = pointer.From(viper.GetBool(fn))
			}

			z, _, err := client.Must().DnsClient.ZonesApi.ZonesPut(context.Background(), uuidgen.Must()).
				ZoneEnsure(dns.ZoneEnsure{Properties: input}).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(z)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Domain name this zone is authoritative for, e.g. example.com", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Free-text note for your own reference; not served in DNS")
	cmd.AddBoolFlag(constants.FlagEnabled, "", true, "Whether the zone is served. true = IONOS answers lookups; false = zone kept but not resolved (default true)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
