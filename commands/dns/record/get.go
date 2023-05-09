package record

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/viper"
)

func ZonesRecordsFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a record",
		Example:   "ionosctl dns record get --zoneId ZONE_ID --recordId RECORD_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagZoneId)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagRecordId)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			r, _, err := client.Must().DnsClient.RecordsApi.ZonesRecordsFindById(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagZoneId)),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagRecordId)),
			).Execute()

			if err != nil {
				return err
			}
			return c.Printer.Print(getRecordPrint(c, r))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagZoneId, "", "", "The ID (UUID) of the DNS zone of which record belongs to")
	cmd.AddStringFlag(constants.FlagRecordId, "", "", "The ID (UUID) of the DNS record")

	return cmd
}
