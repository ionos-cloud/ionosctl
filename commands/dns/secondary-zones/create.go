package secondary_zones

import (
	"context"
	"fmt"
	"net"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	dns "github.com/ionos-cloud/sdk-go-dns"
)

func createCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "create",
			Aliases:   []string{"c"},
			ShortDesc: "Create a secondary zone",
			LongDesc:  `Create a new secondary zone with default NS and SOA records. Note that Cloud DNS relies on the following Anycast addresses for sending DNS notify messages. Make sure to whitelist on your end:

IPv4: 212.227.123.25
IPv6: 2001:8d8:fe:53::5cd:25`,
			Example: "ionosctl dns secondary-zone create --name ZONE_NAME --description DESCRIPTION --primary-ips 1.2.3.4,5.6.7.8",
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
					ZoneName:    &name,
					Description: &description,
					PrimaryIps:  &primaryIPs,
				}

				secZone, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesPost(context.Background()).SecondaryZoneCreate(
					*dns.NewSecondaryZoneCreate(secZoneProps),
				).Execute()
				if err != nil {
					return err
				}

				cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				out, err := jsontabwriter.GenerateOutput(
					"", jsonpaths.DnsSecondaryZone, secZone, tabheaders.GetHeadersAllDefault(allCols, cols),
				)
				if err != nil {
					return err
				}

				fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
				return nil
			},
		},
	)

	c.Command.Flags().StringP(constants.FlagName, constants.FlagNameShort, "", "Name of the secondary zone")
	c.Command.Flags().String(constants.FlagDescription, "", "Description of the secondary zone")
	c.Command.Flags().StringSlice(constants.FlagPrimaryIPs, []string{}, "Primary DNS server IP addresses")

	c.Command.SilenceUsage = true
	c.Command.Flags().SortFlags = false

	return c
}