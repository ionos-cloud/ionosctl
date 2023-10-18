package groups

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	vmasc "github.com/ionos-cloud/sdk-go-vmautoscaling"
	"github.com/spf13/viper"
)

func GroupCreateCmd() *core.Command {

	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "user",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create MongoDB users.",
		Example: fmt.Sprintf("ionosctl autoscaling group create %s",
			core.FlagsUsage(constants.FlagDatacenterId, constants.FlagName)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			return err
		},
		CmdRun: func(c *core.CommandConfig) error {
			execute, a, err := client.Must().VMAscClient.GroupsPost(context.Background()).GroupPost(vmasc.GroupPost{
				Properties: &vmasc.GroupProperties{
					Datacenter: &vmasc.GroupPropertiesDatacenter{
						Id: pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))),
					},
				},
			}).Execute()
			if err != nil {
				return err
			}

			return nil
		},
	})

	return cmd
}
