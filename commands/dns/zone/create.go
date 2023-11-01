package zone

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	dns "github.com/ionos-cloud/sdk-go-dns"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func ZonesPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Create a zone",
		Example:   "ionosctl dns z create --name name.com",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := dns.Zone{}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.ZoneName = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
				input.Description = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagEnabled); viper.IsSet(fn) {
				input.Enabled = pointer.From(viper.GetBool(fn))
			}

			z, _, err := client.Must().DnsClient.ZonesApi.ZonesPut(context.Background(), uuidgen.Must()).
				ZoneEnsure(dns.ZoneEnsure{Properties: &input}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			// if err != nil {
			//	return err
			// }

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DnsZone, z, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the DNS zone, e.g. foo.com")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "The description of the DNS zone")
	cmd.AddBoolFlag(constants.FlagEnabled, "", true, "Activate or deactivate the DNS zone")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
