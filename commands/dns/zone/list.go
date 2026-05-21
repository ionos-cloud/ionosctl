package zone

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/spf13/cobra"
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
			return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
				dnsClient := dns.NewAPIClient(cfg)
				req := dnsClient.ZonesApi.ZonesGet(context.Background())

				if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
					req = req.FilterZoneName(viper.GetString(fn))
				}
				if fn := core.GetFlagName(c.NS, constants.FlagState); viper.IsSet(fn) {
					req = req.FilterState(dns.ProvisioningState(viper.GetString(fn)))
				}

				ls, _, err := req.Execute()
				return ls, err
			})
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
