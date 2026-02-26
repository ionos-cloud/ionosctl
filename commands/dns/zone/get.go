package zone

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func ZonesFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a zone",
		Example:   "ionosctl dns z get --zone ZONE",
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

			z, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), zoneId).Execute()

			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, z, cols))
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
