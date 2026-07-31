package action

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
		Resource:  "action",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a single VM Auto Scaling action",
		LongDesc:  "Show one scaling action of a group by its ID. Because actions are scoped to a group, both --group-id (which group the action belongs to) and --action-id (the action itself) are required. The result includes the action type (SCALE_IN / SCALE_OUT) and its status (IN_PROGRESS / SUCCESSFUL / FAILED) - useful for confirming whether a specific scaling event completed.",
		Example: fmt.Sprintf("ionosctl vm-autoscaling action get %s",
			core.FlagsUsage(constants.FlagGroupId, constants.FlagActionId)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS,
				constants.FlagGroupId, constants.FlagActionId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, _, err := client.Must().VMAscClient.GroupsActionsFindById(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagGroupId)),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagActionId))).
				Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(ls)
		},
	})

	cmd.AddStringFlag(constants.FlagGroupId, "", "", "The ID of the VM Auto Scaling group the action belongs to")
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
	cmd.AddStringFlag(constants.FlagActionId, constants.FlagIdShort, "", "The ID of the scaling action to show (must belong to --group-id)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return ActionsProperty(func(r vmasc.Action) string {
			return fmt.Sprintf("%s\t%s", *r.Id, string(*r.Properties.ActionType))
		}), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
