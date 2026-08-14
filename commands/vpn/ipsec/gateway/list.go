package gateway

import (
	"context"

	vpn "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

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
			return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
				vpnClient := vpn.NewAPIClient(cfg)
				ls, _, err := vpnClient.IPSecGatewaysApi.IpsecgatewaysGet(context.Background()).Execute()
				return ls, err
			})
		},
		InitClient: true,
	})

	return cmd
}
