package server

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
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
