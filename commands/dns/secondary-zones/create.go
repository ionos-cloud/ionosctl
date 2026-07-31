package secondary_zones

import (
	"context"
	"fmt"
	"net"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
)

func createCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "create",
			Aliases:   []string{"c"},
			ShortDesc: "Create a secondary DNS zone",
			LongDesc: `Create a secondary zone: a read-only copy that IONOS pulls from your external primary name server.

--name is the domain (e.g. example.com); --primary-ips lists the primary servers IONOS transfers from. The zone starts with default NS/SOA records; run 'dns secondary-zone transfer start' to pull the rest.

Cloud DNS sends its DNS NOTIFY messages from these Anycast addresses — allow them on your primary:
  IPv4: 212.227.123.25
  IPv6: 2001:8d8:fe:53::5cd:25`,
			Example: "ionosctl dns secondary-zone create --name example.com --primary-ips 1.2.3.4,5.6.7.8",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagPrimaryIPs); err != nil {
					return err
				}

				// Validate primary IPs
				primaryIPs, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagPrimaryIPs)
				for _, ip := range primaryIPs {
					if net.ParseIP(ip) == nil {
						return fmt.Errorf("invalid IP address: %s", ip)
					}
				}

				return nil
			},
			CmdRun: func(c *core.CommandConfig) error {
				name, _ := c.Command.Command.Flags().GetString(constants.FlagName)
				description, _ := c.Command.Command.Flags().GetString(constants.FlagDescription)
				primaryIPs, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagPrimaryIPs)

				secZoneProps := dns.SecondaryZone{
					ZoneName:    name,
					Description: &description,
					PrimaryIps:  primaryIPs,
				}

				secZone, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesPost(context.Background()).SecondaryZoneCreate(
					*dns.NewSecondaryZoneCreate(secZoneProps),
				).Execute()
				if err != nil {
					return err
				}

				return c.Printer(allCols).Print(secZone)
			},
		},
	)

	c.Command.Flags().StringP(constants.FlagName, constants.FlagNameShort, "", "Domain name of the zone to mirror, e.g. example.com")
	c.Command.Flags().String(constants.FlagDescription, "", "Free-text note for your own reference; not served in DNS")
	c.Command.Flags().StringSlice(constants.FlagPrimaryIPs, []string{}, "Comma-separated IPs of the external primary name servers IONOS transfers the zone from")

	c.Command.SilenceUsage = true
	c.Command.Flags().SortFlags = false

	return c
}
