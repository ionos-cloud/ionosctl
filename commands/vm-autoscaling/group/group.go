package group

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
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

	cmd.AddColsFlag(allCols)

	return cmd
}

var allCols = []table.Column{
	{Name: "GroupId", JSONPath: "id", Default: true},
	{Name: "DatacenterId", JSONPath: "properties.datacenter.id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "MinReplicas", JSONPath: "properties.minReplicaCount", Default: true},
	{Name: "Replicas", Default: true, Format: func(item map[string]any) any {
		ents := table.Navigate(item, "entities.servers.items")
		if ents == nil {
			return nil
		}
		arr, ok := ents.([]any)
		if !ok {
			return nil
		}
		return fmt.Sprintf("%d", len(arr))
	}},
	{Name: "MaxReplicas", JSONPath: "properties.maxReplicaCount", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "Metric", JSONPath: "properties.policy.metric"},
	{Name: "Range", JSONPath: "properties.policy.range"},
	{Name: "ScaleInActionAmount", JSONPath: "properties.policy.scaleInAction.amount"},
	{Name: "ScaleInActionAmountType", JSONPath: "properties.policy.scaleInAction.amountType"},
	{Name: "ScaleInActionCooldownPeriod", JSONPath: "properties.policy.scaleInAction.cooldownPeriod"},
	{Name: "ScaleInActionTerminationPolicy", JSONPath: "properties.policy.scaleInAction.terminationPolicy"},
	{Name: "ScaleInActionDeleteVolumes", JSONPath: "properties.policy.scaleInAction.deleteVolumes"},
	{Name: "ScaleInThreshold", JSONPath: "properties.policy.scaleInThreshold"},
	{Name: "ScaleOutActionAmount", JSONPath: "properties.policy.scaleOutAction.amount"},
	{Name: "ScaleOutActionAmountType", JSONPath: "properties.policy.scaleOutAction.amountType"},
	{Name: "ScaleOutActionCooldownPeriod", JSONPath: "properties.policy.scaleOutAction.cooldownPeriod"},
	{Name: "ScaleOutThreshold", JSONPath: "properties.policy.scaleOutThreshold"},
	{Name: "Unit", JSONPath: "properties.policy.unit"},
	{Name: "AvailabilityZone", JSONPath: "properties.replicaConfiguration.availabilityZone"},
	{Name: "Cores", JSONPath: "properties.replicaConfiguration.cores"},
	{Name: "CPUFamily", JSONPath: "properties.replicaConfiguration.cpuFamily"},
	{Name: "RAM", JSONPath: "properties.replicaConfiguration.ram"},
}

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
