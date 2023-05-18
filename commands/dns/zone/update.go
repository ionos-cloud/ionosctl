package zone

import (
	"context"

	dns "github.com/ionos-cloud/sdk-go-dnsaas"

	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/viper"
)

var id string

func ZonesPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a zone's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		Example:   "ionosctl dns zone update --zone-id ZONE_ID --name newname.com",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagZoneId)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			z, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), id).Execute()
			if err != nil {
				return err
			}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				z.Properties.ZoneName = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
				z.Properties.Description = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagEnabled); viper.IsSet(fn) {
				z.Properties.Enabled = pointer.From(viper.GetBool(fn))
			}

			zNew, _, err := client.Must().DnsClient.ZonesApi.ZonesPut(context.Background(), id).
				ZoneUpdateRequest(
					dns.ZoneUpdateRequest{Properties: &dns.ZoneUpdateRequestProperties{
						// We can't pass `z.Properties` directly as it is a different object type
						ZoneName:    z.Properties.ZoneName,
						Description: z.Properties.Description,
						Enabled:     z.Properties.Enabled,
					}},
				).Execute()
			if err != nil {
				return err
			}
			return c.Printer.Print(getZonePrint(c, zNew))
		},
		InitClient: true,
	})

	cmd.AddStringVarFlag(&id, constants.FlagZoneId, constants.FlagIdShort, "", "The ID (UUID) of the DNS zone", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZoneId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return Zones(func(t dns.ZoneResponse) string {
			return *t.GetId()
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the DNS zone, e.g. foo.com")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "The description of the DNS zone")
	cmd.AddBoolFlag(constants.FlagEnabled, "", true, "Activate or deactivate the DNS zone")

	return cmd
}
