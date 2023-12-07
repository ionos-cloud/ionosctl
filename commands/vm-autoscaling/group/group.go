package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "group",
			Aliases:          []string{"g", "groups"},
			Short:            "Autoscaling Groups Operations",
			Long:             "The sub-commands of `ionosctl autoscaling group` allow you to manage the Autoscaling Groups under your account.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(Create())
	cmd.AddCommand(Put())
	cmd.AddCommand(List())
	cmd.AddCommand(Get())
	cmd.AddCommand(Delete())

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

var (
	allCols = []string{
		"GroupId", "DatacenterId", "Name", "MinReplicas", "Replicas", "MaxReplicas", "Location", "State",
		"Metric", "Range", "ScaleInActionAmount", "ScaleInActionAmountType",
		"ScaleInActionCooldownPeriod", "ScaleInActionTerminationPolicy", "ScaleInActionDeleteVolumes",
		"ScaleInThreshold", "ScaleOutActionAmount", "ScaleOutActionAmountType",
		"ScaleOutActionCooldownPeriod", "ScaleOutThreshold", "Unit",
		"AvailabilityZone", "Cores", "CPUFamily", "RAM",
	}

	defaultCols = allCols[0:8]
)

// Groups returns all groups matching the given filters
func Groups(fs ...Filter) (vmasc.GroupCollection, error) {
	req := client.Must().VMAscClient.GroupsGet(context.Background())

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return vmasc.GroupCollection{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return vmasc.GroupCollection{}, err
	}
	return ls, nil
}

func GroupsProperty[V any](f func(vmasc.Group) V, fs ...Filter) []V {
	recs, err := Groups(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

type Filter func(request vmasc.ApiGroupsGetRequest) (vmasc.ApiGroupsGetRequest, error)
