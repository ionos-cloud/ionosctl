package action

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/group"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "action",
			Aliases:          []string{"act"},
			Short:            "Autoscaling Actions Operations",
			Long:             "The sub-commands of `ionosctl autoscaling action` allow you to manage the Autoscaling Actions under your account.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(List())
	cmd.AddCommand(Get())

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

var (
	allJSONPaths = map[string]string{
		"ActionId": "id",
	}

	allCols = []string{
		"ActionId",
	}

	defaultCols = allCols
)

func Actions(fs ...Filter) (vmasc.ActionCollection, error) {
	groupIds := group.GroupsProperty(func(r vmasc.Group) string {
		if r.Id == nil {
			return ""
		}
		return *r.Id
	})

	// for each group, get actions
	var allActions vmasc.ActionCollection
	allActions.Items = pointer.From(make([]vmasc.Action, 0))
	for _, groupId := range groupIds {
		actions, err := GroupActions(groupId, fs...)
		if err != nil {
			return vmasc.ActionCollection{}, err
		}
		allActions.Items = pointer.From(append(*allActions.Items, *actions.Items...))
	}

	return allActions, nil
}

// GroupActions returns all actions matching the given filters from a specific group
func GroupActions(groupId string, fs ...Filter) (vmasc.ActionCollection, error) {
	req := client.Must().VMAscClient.GroupsActionsGet(context.Background(), groupId)

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return vmasc.ActionCollection{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return vmasc.ActionCollection{}, err
	}
	return ls, nil
}

func ActionsProperty[V any](f func(resource vmasc.Action) V, fs ...Filter) []V {
	recs, err := Actions(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

type Filter func(request vmasc.ApiGroupsActionsGetRequest) (vmasc.ApiGroupsActionsGetRequest, error)
