package gateway

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard gateway",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a WireGuard Gateway",
		Example:   "ionosctl ",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			var input *vpn.WireGuardGateway
			if err := setPropertiesFromFlags(c, input); err != nil {
				return err
			}

			id := uuidgen.Must()
			res, _, err := client.Must().CDNClient.DistributionsApi.DistributionsPut(context.Background(), id).
				DistributionUpdate(cdn.DistributionUpdate{
					Id:         &id,
					Properties: input,
				}).Execute()
			if err != nil {
				return err
			}

			return printDistribution(c, res)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return addDistributionCreateFlags(cmd)
}
