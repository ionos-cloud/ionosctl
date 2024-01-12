package jsonpaths

var VmAutoscalingGroup = map[string]string{
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