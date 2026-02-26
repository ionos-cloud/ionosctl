package secondary_zones

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
)

func listCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "list",
			ShortDesc: "List secondary zones",
			LongDesc:  "List all secondary zones. Default limit is the first 100 items. Use pagination query parameters for listing more items (up to 1000).",
			Example:   "ionosctl dns secondary-zone list",
			PreCmdRun: nil,
			CmdRun: func(c *core.CommandConfig) error {
				req := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesGet(context.Background())

				if c.Command.Command.Flags().Changed(constants.FlagName) {
					name, _ := c.Command.Command.Flags().GetString(constants.FlagName)
					req = req.FilterZoneName(name)
				}

				if c.Command.Command.Flags().Changed(constants.FlagState) {
					state, _ := c.Command.Command.Flags().GetString(constants.FlagState)
					req = req.FilterState(dns.ProvisioningState(state))
				}

				secZones, _, err := req.Execute()
				if err != nil {
					return err
				}

				cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				return c.Out(table.Sprint(allCols, secZones, cols, table.WithPrefix("items")))
			},
		},
	)

	enumStates := []string{"AVAILABLE", "FAILED", "PROVISIONING", "DESTROYING"}
	c.AddStringFlag(constants.FlagState, "", "", fmt.Sprintf("Filter used to fetch all zones in a particular state (%s)", strings.Join(enumStates, ", ")))
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagState, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return enumStates, cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.Command.Flags().StringP(constants.FlagName, constants.FlagNameShort, "", "Filter used to fetch only the zones that contain the specified zone name")

	c.Command.SilenceUsage = true
	c.Command.Flags().SortFlags = false

	return c
}
