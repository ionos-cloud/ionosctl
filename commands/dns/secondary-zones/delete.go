package secondary_zones

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/spf13/viper"
)

func deleteCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "delete",
			Aliases:   []string{"d", "del"},
			ShortDesc: "Delete a secondary zone",
			LongDesc:  "Delete a secondary zone",
			Example:   "ionosctl dns secondary-zone delete --zone ZONE_ID",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagZone})
			},
			CmdRun: func(c *core.CommandConfig) error {
				if all, _ := c.Command.Command.Flags().GetBool(constants.ArgAll); c.Command.Command.Flags().Changed(constants.ArgAll) && all {
					return deleteAll(c)
				}

				zone, _ := c.Command.Command.Flags().GetString(constants.FlagZone)
				zoneID, err := utils.SecondaryZoneResolve(zone)
				if err != nil {
					return err
				}

				if err = deleteSingle(c, zoneID); err != nil {
					return err
				}

				c.Msg("Successfully deleted secondary zone %v", zoneID)
				return nil
			},
		},
	)

	c.Command.Flags().BoolP(constants.ArgAll, constants.ArgAllShort, false, "Delete all secondary zones")

	c.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone,
		core.WithCompletion(completer.SecondaryZonesIDs, constants.DNSApiRegionalURL, constants.DNSLocations),
	)

	c.Command.SilenceUsage = true
	c.Command.Flags().SortFlags = false

	return c
}

func deleteAll(c *core.CommandConfig) error {
	return core.DeleteAll(c, core.DeleteAllOptions[dns.SecondaryZoneRead]{
		Resource: "secondary zone",
		List: func() ([]dns.SecondaryZoneRead, error) {
			secZones, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesGet(context.Background()).Execute()
			if err != nil {
				return nil, err
			}
			return secZones.Items, nil
		},
		Summary: func(item dns.SecondaryZoneRead) string {
			s := fmt.Sprintf("%s (id: %s", item.Properties.ZoneName, item.Id)
			if item.Properties.Description != nil && *item.Properties.Description != "" {
				s += fmt.Sprintf(", desc: %s", *item.Properties.Description)
			}
			return s + ")"
		},
		ID: func(item dns.SecondaryZoneRead) string { return item.Id },
		Delete: func(item dns.SecondaryZoneRead) error {
			_, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesDelete(context.Background(), item.Id).Execute()
			return err
		},
	})
}

func deleteSingle(c *core.CommandConfig, zoneId string) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete secondary zone %s", zoneId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesDelete(context.Background(), zoneId).Execute()
	if err != nil {
		return err
	}

	return nil
}
