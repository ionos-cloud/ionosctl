package file

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
)

func getCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "get",
			Aliases:   []string{"g"},
			ShortDesc: "Get a specific zone file",
			LongDesc:  "",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone)
			},
			CmdRun: func(c *core.CommandConfig) error {
				zoneID, _ := c.Command.Command.Flags().GetString(constants.FlagZone)

				resp, err := client.Must().DnsClient.ZoneFilesApi.ZonesZonefileGet(context.Background(), zoneID).Execute()
				if err != nil {
					return err
				}

				fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(string(resp.Payload)))
				return nil
			},
		},
	)

	c.Command.Flags().StringP(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone)

	return c
}
