package gateway

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "ipsec gateway",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List IPSec Gateways",
		Example:   "ionosctl vpn ipsec gateway list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, err := completer.Gateways()
			if err != nil {
				return fmt.Errorf("failed listing gateways: %w", err)
			}

			return c.Printer(allCols).Prefix("items").Print(ls)
		},
	})

	return cmd
}
