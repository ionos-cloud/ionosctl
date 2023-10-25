package groups

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	vmasc "github.com/ionos-cloud/sdk-go-vmautoscaling"
	"github.com/spf13/viper"
)

func GroupCreateCmd() *core.Command {
	var groupProperties vmasc.GroupProperties
	cmd := core.NewCommandWithJsonProperties(context.Background(), nil, string(exampleJson), &groupProperties, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "groups",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create VM Autoscaling Groups",
		Example: fmt.Sprintf("ionosctl vm-autoscaling group create %s",
			core.FlagsUsage(constants.FlagDatacenterId, constants.FlagName)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			dbg, err := json.MarshalIndent(groupProperties, "", "  ")
			if err != nil {
				return fmt.Errorf("failed marshalling viper cfg into group properties: %w", err)
			}
			fmt.Println(string(dbg))

			group, _, err := client.Must().VMAscClient.GroupsPost(context.Background()).GroupPost(vmasc.GroupPost{
				Properties: &groupProperties,
			}).Execute()
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

	return cmd
}

var exampleJson, _ = json.MarshalIndent(`
	{
"properties": {
"datacenter": {
"id": "09043280-aafc-49f7-a048-d61673f52024"
},
"maxReplicaCount": 10,
"minReplicaCount": 1,
"name": "VM Auto Scaling Group 1",
"policy": {
"metric": "INSTANCE_CPU_UTILIZATION_AVERAGE",
"range": "P1D",
"scaleInAction": {
"amount": 1,
"amountType": "ABSOLUTE",
"cooldownPeriod": "5m",
"terminationPolicy": "OLDEST_SERVER_FIRST",
"deleteVolumes": true
},
"scaleInThreshold": 33,
"scaleOutAction": {
"amount": 1,
"amountType": "ABSOLUTE",
"cooldownPeriod": "5m"
},
"scaleOutThreshold": 77,
"unit": "PER_MINUTE"
},
"replicaConfiguration": {
"availabilityZone": "AUTO",
"cores": 2,
"cpuFamily": "INTEL_SKYLAKE",
"nics": [
{
"lan": 1,
"name": "LAN NIC 1",
"dhcp": true,
"firewallActive": true,
"firewallType": "INGRESS",
"flowLogs": [
{
"name": "flow-log",
"action": "ACCEPTED",
"direction": "EGRESS",
"bucket": "bucketName/key"
}
],
"firewallRules": [
{
"name": "My resource",
"protocol": "TCP",
"sourceMac": "00:0a:95:9d:68:16",
"sourceIp": "22.231.113.64",
"targetIp": "22.231.113.64",
"icmpCode": 2,
"icmpType": 8,
"portRangeStart": 8,
"portRangeEnd": 8,
"type": "INGRESS"
}
],
"targetGroup": {
"targetGroupId": "id_example",
"port": 8080,
"weight": 15
}
}
],
"ram": 2048,
"volumes": [
{
"image": "6e928bd0-3a8e-4821-a20a-54984b0c2d21",
"imageAlias": "ubuntu:latest",
"name": "Volume 1",
"size": 30,
"sshKeys": [
"ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSU\nGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3\nPbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XA\nt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/En\nmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbx\nNrRFi9wrf+M7Q== user@domain.local"
],
"type": "SSD",
"userData": "ZWNobyAiSGVsbG8sIFdvcmxkIgo=",
"bus": "VIRTIO",
"backupunitId": "25f67991-0f51-4efc-a8ad-ef1fb31a481c",
"bootOrder": "AUTO",
"imagePassword": "passw0rd"
}
]
}
}
}
`, "", "  ")
