package zone

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func ZonesPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a zone's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		Example:   "ionosctl dns z update --zone ZONE --name newname.com",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagZone); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			zone, _ := c.Command.Command.Flags().GetString(constants.FlagZone)
			id, err := utils.ZoneResolve(zone)
			if err != nil {
				return err
			}

			z, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), id).Execute()
			if err != nil {
				return err
			}

			if c.Command.Command.Flags().Changed(constants.FlagName) {
				name, _ := c.Command.Command.Flags().GetString(constants.FlagName)
				z.Properties.ZoneName = name
			}
			if c.Command.Command.Flags().Changed(constants.FlagDescription) {
				desc, _ := c.Command.Command.Flags().GetString(constants.FlagDescription)
				z.Properties.Description = pointer.From(desc)
			}
			if c.Command.Command.Flags().Changed(constants.FlagEnabled) {
				enabled, _ := c.Command.Command.Flags().GetBool(constants.FlagEnabled)
				z.Properties.Enabled = pointer.From(enabled)
			}

			zNew, _, err := client.Must().DnsClient.ZonesApi.ZonesPut(context.Background(), id).
				ZoneEnsure(dns.ZoneEnsure{Properties: z.Properties}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DnsZone, zNew, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
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
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The new name of the DNS zone, e.g. foo.com")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "The new description of the DNS zone")
	cmd.AddBoolFlag(constants.FlagEnabled, "", true, "Activate or deactivate the DNS zone")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
