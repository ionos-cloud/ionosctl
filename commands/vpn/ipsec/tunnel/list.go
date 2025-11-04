package tunnel

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "ipsec tunnel",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List IPSec Tunnels",
		Example:   "ionosctl vpn ipsec tunnel list " + core.FlagsUsage(constants.FlagGatewayID),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID)
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, err := completer.Tunnels(viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID)))
			if err != nil {
				return fmt.Errorf("failed listing tunnels: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("items", jsonpaths.VPNIPSecTunnel, ls,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the IPSec Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL, constants.VPNLocations),
	)

	return cmd
}
