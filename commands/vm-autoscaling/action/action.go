package action

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/group"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "action",
			Aliases: []string{"act"},
			Short:   "Inspect the scaling history of a VM Auto Scaling group",
			Long: `A VM Auto Scaling action is a single scaling event that the group performed: a SCALE_OUT that added replicas or a SCALE_IN that removed them, triggered when the group's metric crossed a policy threshold. Actions are the audit trail of the autoscaler - you do not create them, the group does.

Each action records its type (SCALE_IN / SCALE_OUT) and a status: IN_PROGRESS (still executing), SUCCESSFUL, or FAILED. Reviewing actions is how you answer "why did my replica count change, and when?" and how you spot scaling that failed (e.g. because a scale-out hit a resource limit).

Actions belong to a group, so list them per group with --group-id, or across all your groups with --all.`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(List())
	cmd.AddCommand(Get())

	cmd.AddColsFlag(allCols)

	return cmd
}

var allCols = []table.Column{
	{Name: "ActionId", JSONPath: "id", Default: true},
	{Name: "GroupId", JSONPath: "href", Default: true},
}

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
