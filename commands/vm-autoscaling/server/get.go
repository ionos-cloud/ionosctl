package server

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/group"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "server",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a VM Autoscaling Server",
		Example: fmt.Sprintf("ionosctl vm-autoscaling server get %s",
			core.FlagsUsage(constants.FlagGroupId, constants.FlagServerId)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS,
				constants.FlagGroupId, constants.FlagServerId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, _, err := client.Must().VMAscClient.GroupsServersFindById(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagGroupId)),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagServerId))).
				Execute()
			if err != nil {
				return err
			}

			enriched, err := enrichAutoscalingServer(ls)
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(enriched)
		},
	})

	cmd.AddStringFlag(constants.FlagGroupId, "", "", "ID of the autoscaling group that the server is a part of")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return group.GroupsProperty(func(r vmasc.Group) string {
			completion := *r.Id
			if r.Properties == nil || r.Properties.Name == nil {
				return completion
			}
			completion += "\t" + *r.Properties.Name
			return completion
		}, func(r vmasc.ApiGroupsGetRequest) (vmasc.ApiGroupsGetRequest, error) {
			return r.Depth(1), nil
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagServerId, constants.FlagIdShort, "", "ID of the autoscaling server")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return ServersProperty(func(r vmasc.Server) string {
			return fmt.Sprintf("%s\t%s", *r.Id, *r.Properties.Name)
		}), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
