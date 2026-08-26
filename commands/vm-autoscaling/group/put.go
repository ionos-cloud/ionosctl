package group

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Put() *core.Command {
	var jsonStruct vmasc.GroupPut
	cmd := core.NewCommandWithJsonProperties(context.Background(), nil, exampleJson, &jsonStruct, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "groups",
		Verb:      "put",
		Aliases:   []string{"p", "update"},
		ShortDesc: "Replace the configuration of an existing VM Auto Scaling group",
		LongDesc: `Replace the properties of an existing VM Auto Scaling group (HTTP PUT). This is a full replacement, not a patch: the JSON 'properties' object you pass becomes the group's new configuration in its entirety, so include every field you want to keep. A common workflow is to 'group get' the current group, edit the properties, and pass them back here.

Typical reasons to update a group:
  * Change the replica bounds (minReplicaCount / maxReplicaCount) to widen or narrow how far it can scale. Lowering maxReplicaCount below the current replica count triggers scale-in; raising minReplicaCount above it triggers scale-out.
  * Retune the policy - e.g. adjust scaleInThreshold / scaleOutThreshold, switch the metric or unit, or change the scale action amount and cooldownPeriod.
  * Update the replicaConfiguration template (image, cores, ram, NICs). Note this affects replicas created AFTER the change; existing replicas are not re-provisioned.

See the field-by-field reference under 'group create' for the meaning and valid values of each property.`,
		Example: fmt.Sprintf(`# Widen an existing group's replica bounds and raise its scale-out threshold
ionosctl vm-autoscaling group put %s --json-properties '{"properties":{"datacenter":{"id":"<datacenter-id>"},"name":"web-tier","minReplicaCount":2,"maxReplicaCount":20,"policy":{"metric":"INSTANCE_CPU_UTILIZATION_AVERAGE","range":"2m","unit":"PER_MINUTE","scaleInThreshold":30,"scaleOutThreshold":85,"scaleInAction":{"amount":1,"amountType":"ABSOLUTE","cooldownPeriod":"5m","terminationPolicy":"OLDEST_SERVER_FIRST","deleteVolumes":true},"scaleOutAction":{"amount":2,"amountType":"ABSOLUTE","cooldownPeriod":"10m"}},"replicaConfiguration":{"availabilityZone":"AUTO","cores":2,"cpuFamily":"INTEL_SKYLAKE","ram":2048,"nics":[{"lan":1,"name":"nic1","dhcp":true}],"volumes":[{"imageAlias":"ubuntu:latest","name":"boot","size":30,"type":"SSD","imagePassword":"<password>"}]}}}'`,
			core.FlagUsage(constants.FlagGroupId)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagGroupId, constants.FlagJsonProperties},
				[]string{constants.FlagGroupId, constants.FlagJsonPropertiesExample},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			group, _, err := client.Must().VMAscClient.GroupsPut(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagGroupId))).
				GroupPut(jsonStruct).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(group)
		},
	})

	cmd.AddStringFlag(constants.FlagGroupId, constants.FlagIdShort, "", "The ID of the VM Auto Scaling group whose configuration is replaced", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// get ID of all groups
		return GroupsProperty(func(r vmasc.Group) string {
			completion := *r.Id
			if r.Properties == nil || r.Properties.Name == nil {
				return completion
			}
			completion += "\t" + *r.Properties.Name
			return completion
		}, func(r vmasc.ApiGroupsGetRequest) (vmasc.ApiGroupsGetRequest, error) {
			return r.Depth(1), nil
		}), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
