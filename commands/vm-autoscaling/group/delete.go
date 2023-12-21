package group

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Delete() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "groups",
		Verb:      "delete",
		Aliases:   []string{"d", "del", "rm"},
		ShortDesc: "Delete VM Autoscaling Groups",
		Example: fmt.Sprintf("ionosctl vm-autoscaling group delete (%s|--%s)",
			core.FlagUsage(constants.FlagGroupId), constants.ArgAll),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagGroupId},
				[]string{constants.ArgAll},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
				return deleteGroups(c, GroupsProperty(func(r vmasc.Group) string {
					return *r.Id
				}))
			}
			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagGroupId))
			return deleteGroups(c, []string{id})
		},
	})

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Set this flag to delete all VM-Autoscaling groups from your account")
	cmd.AddStringFlag(constants.FlagGroupId, constants.FlagIdShort, "", "ID of the autoscaling group to list servers from")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// get ID of all groups
		return GroupsProperty(func(r vmasc.Group) string {
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
func deleteGroups(c *core.CommandConfig, ids []string) error {
	var errs error
	for _, id := range ids {
		group, _, err := client.Must().VMAscClient.GroupsFindById(context.Background(), id).Execute()
		if err != nil {
			return fmt.Errorf("failed retrieving info about group %s: %w", id, err)
		}

		if shouldDeleteGroup(c, &group) {
			_, err := client.Must().VMAscClient.GroupsDelete(context.Background(), id).Execute()
			if err != nil {
				return fmt.Errorf("failed deleting group %s: %w", id, err)
			}
		} else {
			errs = errors.Join(errs, fmt.Errorf("%s for %s", confirm.UserDenied, *group.Id))
		}
	}

	return errs
}

func shouldDeleteGroup(c *core.CommandConfig, group *vmasc.Group) bool {
	return confirm.FAsk(c.Command.Command.InOrStdin(),
		fmt.Sprintf("Do you really want to delete group %s from %s (%s)", *group.Properties.Name, *group.Properties.Location, *group.Id),
		viper.GetBool(constants.ArgForce))
}
