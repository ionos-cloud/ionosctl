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

	cmd.AddCommand(GroupCreateCmd())
	cmd.AddCommand(GroupListCmd())
	cmd.AddCommand(GroupDeleteCmd())

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

var (
	allJSONPaths = map[string]string{
		"GroupId":                        "id",
		"Name":                           "properties.name",
		"MinReplicas":                    "properties.minReplicaCount",
		"MaxReplicas":                    "properties.maxReplicaCount",
		"DatacenterId":                   "properties.datacenter.id",
		"State":                          "metadata.state",
		"Location":                       "properties.location",
		"Metric":                         "properties.policy.metric",
		"Range":                          "properties.policy.range",
		"ScaleInActionAmount":            "properties.policy.scaleInAction.amount",
		"ScaleInActionAmountType":        "properties.policy.scaleInAction.amountType",
		"ScaleInActionCooldownPeriod":    "properties.policy.scaleInAction.cooldownPeriod",
		"ScaleInActionTerminationPolicy": "properties.policy.scaleInAction.terminationPolicy",
		"ScaleInActionDeleteVolumes":     "properties.policy.scaleInAction.deleteVolumes",
		"ScaleInThreshold":               "properties.policy.scaleInThreshold",
		"ScaleOutActionAmount":           "properties.policy.scaleOutAction.amount",
		"ScaleOutActionAmountType":       "properties.policy.scaleOutAction.amountType",
		"ScaleOutActionCooldownPeriod":   "properties.policy.scaleOutAction.cooldownPeriod",
		"ScaleOutThreshold":              "properties.policy.scaleOutThreshold",
		"Unit":                           "properties.policy.unit",
		"AvailabilityZone":               "properties.replicaConfiguration.availabilityZone",
		"Cores":                          "properties.replicaConfiguration.cores",
		"CPUFamily":                      "properties.replicaConfiguration.cpuFamily",
		"RAM":                            "properties.replicaConfiguration.ram",
	}

	allCols = []string{
		"GroupId", "Name", "MinReplicas", "MaxReplicas", "DatacenterId", "Location", "State",
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

func GroupsProperty[V any](f func(vmasc.GroupResource) V, fs ...Filter) []V {
	recs, err := Groups(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

type Filter func(request vmasc.ApiGroupsGetRequest) (vmasc.ApiGroupsGetRequest, error)
