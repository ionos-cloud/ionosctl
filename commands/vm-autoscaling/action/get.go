package action

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/group"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "action",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a VM Autoscaling Action",
		Example: fmt.Sprintf("ionosctl vm-autoscaling action get %s",
			core.FlagsUsage(constants.FlagGroupId, constants.FlagActionId)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS,
				constants.FlagGroupId, constants.FlagActionId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			groupId, _ := c.Command.Command.Flags().GetString(constants.FlagGroupId)
			actionId, _ := c.Command.Command.Flags().GetString(constants.FlagActionId)

			ls, _, err := client.Must().VMAscClient.GroupsActionsFindById(context.Background(),
				groupId,
				actionId).
				Execute()
			if err != nil {
				return err
			}

			colsDesired, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("", allJSONPaths, ls,
				tabheaders.GetHeaders(allCols, defaultCols, colsDesired))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagGroupId, "", "", "ID of the autoscaling group that the action is a part of")
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
	cmd.AddStringFlag(constants.FlagActionId, constants.FlagIdShort, "", "ID of the autoscaling action")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return ActionsProperty(func(r vmasc.Action) string {
			return fmt.Sprintf("%s\t%s", *r.Id, string(*r.Properties.ActionType))
		}), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
