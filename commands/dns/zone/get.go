package zone

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/viper"
)

func ZonesFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a zone",
		Example:   "ionosctl dns zone get --zone-id ZONE_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagZoneId)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			z, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagZoneId)),
			).Execute()

			if err != nil {
				return err
			}
			return c.Printer.Print(getZonePrint(c, z))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZoneId, constants.FlagIdShort, "", "The ID (UUID) of the DNS zone", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagZoneId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.Zones(), cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}
