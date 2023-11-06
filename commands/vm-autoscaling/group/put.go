package group

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
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
		Aliases:   []string{"create", "c"},
		ShortDesc: "Guarantee existance of a VM Autoscaling Group. Can be used to either update or create a group.",
		Example: fmt.Sprintf("ionosctl vm-autoscaling group put %s",
			core.FlagsUsage(constants.FlagJsonProperties)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagJsonProperties},
				[]string{constants.FlagJsonPropertiesExample},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			wantedId := uuidgen.Must()
			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagGroupId)) {
				wantedId = viper.GetString(core.GetFlagName(c.NS, constants.FlagGroupId))
			}
			group, _, err := client.Must().VMAscClient.GroupsPut(context.Background(), wantedId).
				GroupPut(jsonStruct).Execute()
			if err != nil {
				return err
			}

			colsDesired := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))
			out, err := jsontabwriter.GenerateOutput("", allJSONPaths, group,
				tabheaders.GetHeaders(allCols, defaultCols, colsDesired))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagGroupId, constants.FlagIdShort, "", "ID of the autoscaling group to create/modify. If not set, a random UUIDv5 will be generated.")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// get ID of all groups
		return GroupsProperty(func(r vmasc.GroupResource) string {
			return fmt.Sprintf(*r.Id) // + "\t" + *r.Properties.Name) // Commented because this SDK functionality currently broken
		}), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

var exampleJson = "{\n  \"properties\": {\n    \"datacenter\": {\n      \"id\": \"6e928bd0-3a8e-4821-a20a-54984b0c2d21\"\n    },\n    \"maxReplicaCount\": 10,\n    \"minReplicaCount\": 1,\n    \"name\": \"AutoScaling-Group1\",\n    \"policy\": {\n      \"metric\": \"INSTANCE_CPU_UTILIZATION_AVERAGE\",\n      \"range\": \"P1D\",\n      \"scaleInAction\": {\n        \"amount\": 1,\n        \"amountType\": \"ABSOLUTE\",\n        \"cooldownPeriod\": \"5m\",\n        \"terminationPolicy\": \"OLDEST_SERVER_FIRST\",\n        \"deleteVolumes\": true\n      },\n      \"scaleInThreshold\": 33,\n      \"scaleOutAction\": {\n        \"amount\": 1,\n        \"amountType\": \"ABSOLUTE\",\n        \"cooldownPeriod\": \"5m\"\n      },\n      \"scaleOutThreshold\": 77,\n      \"unit\": \"PER_MINUTE\"\n    },\n    \"replicaConfiguration\": {\n      \"availabilityZone\": \"AUTO\",\n      \"cores\": 2,\n      \"cpuFamily\": \"INTEL_SKYLAKE\",\n      \"nics\": [\n        {\n          \"lan\": 1,\n          \"name\": \"LAN-NIC-1\",\n          \"dhcp\": true,\n          \"firewallActive\": true,\n          \"firewallType\": \"INGRESS\",\n          \"flowLogs\": [\n            {\n              \"name\": \"flow-log\",\n              \"action\": \"ACCEPTED\",\n              \"direction\": \"EGRESS\",\n              \"bucket\": \"bucketName/key\"\n            }\n          ],\n          \"firewallRules\": [\n            {\n              \"name\": \"My-resource\",\n              \"protocol\": \"TCP\",\n              \"sourceMac\": \"00:0a:95:9d:68:16\",\n              \"sourceIp\": \"22.231.113.64\",\n              \"targetIp\": \"22.231.113.64\",\n              \"icmpCode\": 2,\n              \"icmpType\": 8,\n              \"portRangeStart\": 8,\n              \"portRangeEnd\": 8,\n              \"type\": \"INGRESS\"\n            }\n          ],\n          \"targetGroup\": {\n            \"targetGroupId\": \"6e928bd0-3a8e-4821-a20a-54984b0c2d21\",\n            \"port\": 8080,\n            \"weight\": 15\n          }\n        }\n      ],\n      \"ram\": 2048,\n      \"volumes\": [\n        {\n          \"image\": \"6e928bd0-3a8e-4821-a20a-54984b0c2d21\",\n          \"imageAlias\": \"ubuntu:latest\",\n          \"name\": \"Volume-1\",\n          \"size\": 30,\n          \"sshKeys\": [\n            \"ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSU\\nGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3\\nPbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XA\\nt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/En\\nmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbx\\nNrRFi9wrf+M7Q== user@domain.local\"\n          ],\n          \"type\": \"SSD\",\n          \"userData\": \"ZWNobyAiSGVsbG8sIFdvcmxkIgo=\",\n          \"bus\": \"VIRTIO\",\n          \"backupunitId\": \"25f67991-0f51-4efc-a8ad-ef1fb31a481c\",\n          \"bootOrder\": \"AUTO\",\n          \"imagePassword\": \"passw0rd\"\n        }\n      ]\n    }\n  }\n}\n"
