package secondary_zones

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/cobra"
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
				return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
					dnsClient := dns.NewAPIClient(cfg)
					req := dnsClient.SecondaryZonesApi.SecondaryzonesGet(context.Background())

					if c.Command.Command.Flags().Changed(constants.FlagName) {
						name, _ := c.Command.Command.Flags().GetString(constants.FlagName)
						req = req.FilterZoneName(name)
					}
					if c.Command.Command.Flags().Changed(constants.FlagState) {
						state, _ := c.Command.Command.Flags().GetString(constants.FlagState)
						req = req.FilterState(dns.ProvisioningState(state))
					}

					ls, _, err := req.Execute()
					return ls, err
				})
			},
			InitClient: true,
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
