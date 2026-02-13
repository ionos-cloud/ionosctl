package zone

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func ZonesGetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "Retrieve zones",
		Example:   "ionosctl dns z list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			req := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background())

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				req = req.FilterZoneName(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagState); viper.IsSet(fn) {
				req = req.FilterState(dns.ProvisioningState(viper.GetString(fn)))
			}

			ls, _, err := req.Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("items", jsonpaths.DnsZone, ls,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	enumStates := []string{"AVAILABLE", "FAILED", "PROVISIONING", "DESTROYING"}
	cmd.AddStringFlag(constants.FlagState, "", "", fmt.Sprintf("Filter used to fetch all zones in a particular state (%s)", strings.Join(enumStates, ", ")))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagState, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return enumStates, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Filter used to fetch only the zones that contain the specified zone name")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
