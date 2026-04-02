package server

import (
	"context"

	computeserver "github.com/ionos-cloud/ionosctl/v6/commands/compute/server"
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

var (
	allCols = append([]string{"GroupServerId"}, computeserver.AllServerCols...)

	defaultCols = allCols
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
	cmd.AddCommand(Get())

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(defaultCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func Servers(fs ...Filter) (vmasc.ServerCollection, error) {
	groupIds := group.GroupsProperty(func(r vmasc.Group) string {
		if r.Id == nil {
			return ""
		}
		return *r.Id
	})

	var ls vmasc.ServerCollection
	ls.Items = pointer.From(make([]vmasc.Server, 0))
	for _, groupId := range groupIds {
		actions, err := GroupServers(groupId, fs...)
		if err != nil {
			return vmasc.ServerCollection{}, err
		}
		ls.Items = pointer.From(append(*ls.Items, *actions.Items...))
	}

	return ls, nil
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

func ServersProperty[V any](f func(resource vmasc.Server) V, fs ...Filter) []V {
	recs, err := Servers(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

type Filter func(request vmasc.ApiGroupsServersGetRequest) (vmasc.ApiGroupsServersGetRequest, error)
