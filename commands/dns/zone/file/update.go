package file

import (
	"context"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	dns "github.com/ionos-cloud/sdk-go-dns"
)

func updateCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "update",
			ShortDesc: "Update a zone file",
			LongDesc:  "Update a zone file",
			Example:   "ionosctl dns zone file update --zone ZONE_ID --zone-file FILE_PATH",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone, constants.FlagZoneFile)
			},
			CmdRun: func(c *core.CommandConfig) error {
				zoneNameOrID, _ := c.Command.Command.Flags().GetString(constants.FlagZone)
				zoneID, err := utils.ZoneResolve(zoneNameOrID)
				if err != nil {
					return err
				}

				zoneFilePath, _ := c.Command.Command.Flags().GetString(constants.FlagZoneFile)
				body, err := os.ReadFile(zoneFilePath)
				if err != nil {
					return fmt.Errorf("failed to read zone file: %s", err)
				}

				_, _, err = client.Must().DnsClient.ZoneFilesApi.ZonesZonefilePut(context.Background(), zoneID).Body(string(body)).Execute()
				if err != nil {
					return err
				}

				fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Zone file updated for zone %s", zoneID))
				return nil
			},
		},
	)

	c.AddStringFlag(constants.FlagZone, constants.FlagZoneShort, "", constants.DescZone, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ZonesProperty(func(t dns.ZoneRead) string {
				return *t.Properties.ZoneName
			})
		}, constants.PlaceholderDnsApiURL),
	)
	c.Command.Flags().String(constants.FlagZoneFile, "", "Path to the zone file")

	c.Command.SilenceUsage = true
	c.Command.Flags().SortFlags = false

	return c
}
