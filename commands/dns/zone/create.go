package zone

import (
	"context"

	dns "github.com/ionos-cloud/sdk-go-dnsaas"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Create a zone",
		Example:   "ionosctl dns zone create --name name.com",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagName)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {

			input := dns.ZoneCreateRequestProperties{}
			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.ZoneName = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
				input.Description = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagEnabled); viper.IsSet(fn) {
				input.Enabled = pointer.From(viper.GetBool(fn))
			}

			z, _, err := client.Must().DnsClient.ZonesApi.ZonesPost(context.Background()).
				ZoneCreateRequest(dns.ZoneCreateRequest{Properties: &input}).Execute()
			if err != nil {
				return err
			}
			return c.Printer.Print(getZonePrint(c, z))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the DNS zone, e.g. foo.com")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "The description of the DNS zone")
	cmd.AddBoolFlag(constants.FlagEnabled, "", true, "Activate or deactivate the DNS zone")

	return cmd
}
