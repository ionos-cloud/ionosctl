package group

import (
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/spf13/cobra"
)

func GroupCmd() *core.Command {
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
	cmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

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
