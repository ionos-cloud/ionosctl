package dnssec

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
)

func Delete() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "dnssec",
		Verb:      "delete",
		Aliases:   []string{"del", "rm", "remove"},
		ShortDesc: "Removes ALL associated DNSKEY records for your DNS zone and disables DNSSEC keys.",
		Example:   `ionosctl dns keys delete --zone ZONE`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			zoneName, err := c.Command.Command.Flags().GetString(constants.FlagZone)
			if err != nil {
				return err
			}
			zoneId, err := utils.ZoneResolve(zoneName)
			if err != nil {
				return err
			}

			_, _, err = client.Must().DnsClient.DNSSECApi.ZonesKeysDelete(context.Background(), zoneId).Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "DNSKEY records deleted and DNSSEC keys disabled for zone %s\n", zoneName)
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

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
