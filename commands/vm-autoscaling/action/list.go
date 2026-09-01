package action

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/group"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "action",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List the scaling actions (SCALE_IN / SCALE_OUT history) of one or all groups",
		LongDesc:  "List the scaling actions recorded for a group - the log of every SCALE_IN and SCALE_OUT it performed, each with a status (IN_PROGRESS / SUCCESSFUL / FAILED). Pass --group-id to see the history of one group, or --all to gather actions across every group in your account (this fetches actions group-by-group, so it is slower with many groups).",
		Example: fmt.Sprintf(`ionosctl vm-autoscaling action list %s
ionosctl vm-autoscaling action list %s`,
			core.FlagUsage(constants.FlagGroupId), core.FlagUsage(constants.ArgAll)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagGroupId},
				[]string{constants.ArgAll},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
				return listAll(c)
			}

			// list actions of a group
			ls, err := GroupActions(viper.GetString(core.GetFlagName(c.NS, constants.FlagGroupId)))
			if err != nil {
				return fmt.Errorf("failed listing actions of group %s: %w",
					viper.GetString(core.GetFlagName(c.NS, constants.FlagGroupId)), err)
			}

			return c.Printer(allCols).Prefix("items").Print(ls)
		},
	})

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List actions from every VM Auto Scaling group in your account (mutually exclusive with --group-id)")
	cmd.AddStringFlag(constants.FlagGroupId, constants.FlagIdShort, "", "The ID of the VM Auto Scaling group whose scaling actions to list")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// get ID of all groups
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

	return cmd
}

func listAll(c *core.CommandConfig) error {
	// list actions of all groups
	ls, err := Actions()
	if err != nil {
		return fmt.Errorf("failed listing actions of all groups: %w", err)
	}

	return c.Printer(allCols).Prefix("items").Print(ls)
}
