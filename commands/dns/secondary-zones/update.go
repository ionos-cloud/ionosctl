package secondary_zones

import (
	"context"
	"fmt"
	"net"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"

	dns "github.com/ionos-cloud/sdk-go-dns"
)

func updateCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "update",
			ShortDesc: "Update or create a secondary zone",
			LongDesc:  "Update or create a secondary zone",
			Example:   "ionosctl dns secondary-zone update --zone ZONE_ID --name ZONE_NAME --description DESCRIPTION --primary-ips 1.2.3.4,5.6.7.8",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone); err != nil {
					return err
				}

				if !c.Command.Command.Flags().Changed(constants.FlagDescription) &&
					!c.Command.Command.Flags().Changed(constants.FlagPrimaryIPs) {
					return fmt.Errorf(
						"at least one of the following flags must be set: %s, %s",
						constants.FlagDescription, constants.FlagPrimaryIPs,
					)
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
				zoneNameOrID, _ := c.Command.Command.Flags().GetString(constants.FlagZone)
				zoneID, err := utils.SecondaryZoneResolve(zoneNameOrID)
				if err != nil {
					return err
				}

				secZoneProps, err := setSecondaryZoneProperties(c, zoneID)
				if err != nil {
					return err
				}

				secZone, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesPut(
					context.Background(), zoneID,
				).SecondaryZoneEnsure(*dns.NewSecondaryZoneEnsure(secZoneProps)).Execute()
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

	c.Command.Flags().String(constants.FlagDescription, "", "Description of the secondary zone")
	c.Command.Flags().StringSlice(constants.FlagPrimaryIPs, []string{}, "Primary DNS server IP addresses")
	c.Command.Flags().StringP(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone)

	return c
}

func setSecondaryZoneProperties(c *core.CommandConfig, zoneID string) (dns.SecondaryZone, error) {
	currentZone, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesFindById(context.Background(), zoneID).Execute()
	if err != nil {
		return dns.SecondaryZone{}, err
	}

	if c.Command.Command.Flags().Changed(constants.FlagDescription) {
		description, _ := c.Command.Command.Flags().GetString(constants.FlagDescription)
		currentZone.Properties.Description = &description
	}

	if c.Command.Command.Flags().Changed(constants.FlagPrimaryIPs) {
		primaryIPs, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagPrimaryIPs)
		currentZone.Properties.PrimaryIps = &primaryIPs
	}

	return *currentZone.Properties, nil
}
