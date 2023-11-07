package server

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
			Use:              "server",
			Aliases:          []string{"s", "sv", "vm", "vms", "servers"},
			Short:            "Autoscaling Servers Operations",
			Long:             "The sub-commands of `ionosctl autoscaling server` allow you to manage the Autoscaling Servers under your account.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(List())

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

var (
	allJSONPaths = map[string]string{
		"ServerId": "id",
	}

	allCols = []string{
		"ServerId",
	}

	defaultCols = allCols
)

func Servers(fs ...Filter) (vmasc.ServerCollection, error) {
	groupIds := group.GroupsProperty(func(r vmasc.GroupResource) string {
		if r.Id == nil {
			return ""
		}
		return *r.Id
	})

	// for each group, get actions
	var allActions vmasc.ServerCollection
	allActions.Items = pointer.From(make([]vmasc.ServerResource, 0))
	for _, groupId := range groupIds {
		actions, err := GroupServers(groupId, fs...)
		if err != nil {
			return vmasc.ServerCollection{}, err
		}
		allActions.Items = pointer.From(append(*allActions.Items, *actions.Items...))
	}

	return allActions, nil
}

func GroupServers(groupId string, fs ...Filter) (vmasc.ServerCollection, error) {
	req := client.Must().VMAscClient.GroupsServersGet(context.Background(), groupId)

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return vmasc.ServerCollection{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return vmasc.ServerCollection{}, err
	}
	return ls, nil
}

func ServersProperty[V any](f func(resource vmasc.ServerResource) V, fs ...Filter) []V {
	recs, err := Servers(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

type Filter func(request vmasc.ApiGroupsServersGetRequest) (vmasc.ApiGroupsServersGetRequest, error)
