package zone

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a zone",
		Example:   "ionosctl dns zone delete --zone-id ZONE_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagZoneId)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			_, err := client.Must().DnsClient.ZonesApi.ZonesDelete(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagZoneId)),
			).Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZoneId, constants.FlagIdShort, "", "The ID (UUID) of the DNS zone", core.RequiredFlagOption())

	return cmd
}
