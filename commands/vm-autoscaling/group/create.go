package group

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
)

func Create() *core.Command {
	var jsonStruct vmasc.GroupPost
	cmd := core.NewCommandWithJsonProperties(context.Background(), nil, exampleJson, &jsonStruct, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "groups",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a VM Auto Scaling group from a JSON properties object",
		LongDesc: `Create a new VM Auto Scaling group. The group's full configuration is supplied as a JSON 'properties' object (--json-properties, or a path to a file), because it nests three sub-objects that would be unwieldy as individual flags. Run this command with --json-properties-example to print a ready-to-edit template.

The properties object has these parts:

  * datacenter.id (required) - the datacenter the group and all its replicas live in.
  * name - a friendly label for the group.
  * minReplicaCount / maxReplicaCount (0-100) - the floor and ceiling for the number of running replicas. The autoscaler stays within these bounds.
  * policy - when and how to scale:
      - metric: INSTANCE_CPU_UTILIZATION_AVERAGE | INSTANCE_NETWORK_IN_BYTES | INSTANCE_NETWORK_IN_PACKETS | INSTANCE_NETWORK_OUT_BYTES | INSTANCE_NETWORK_OUT_PACKETS
      - range: the metric aggregation window in Nm/Ns/Nh notation, e.g. 2m or 120s; minimum 120s (API default 120s).
      - unit: PER_SECOND | PER_MINUTE | PER_HOUR | TOTAL - how the metric is normalized. For network metrics this is the rate unit; if unit=TOTAL, scaleOutThreshold must be >= 40.
      - scaleOutThreshold: when the metric goes ABOVE this value, scale out. scaleInThreshold: when it drops BELOW this value, scale in. The two thresholds must keep a minimum gap (metric-dependent) so the group does not scale in and out at the same time. scaleInThreshold < scaleOutThreshold.
      - scaleOutAction / scaleInAction: amount = how many replicas (or what percentage, when amountType=PERCENTAGE) to add/remove; amountType = ABSOLUTE | PERCENTAGE; cooldownPeriod = how long to wait before the next action (min 2m, max 24h, default 5m).
      - scaleInAction.terminationPolicy: OLDEST_SERVER_FIRST | NEWEST_SERVER_FIRST | RANDOM - which replica to remove first. scaleInAction.deleteVolumes: if true, a removed replica's volumes are deleted too. Leave true unless you need to keep the data; orphaned volumes count against your contract limits and can eventually block further scale-outs.
  * replicaConfiguration - the template every new replica is cloned from: cores, ram (in MB), cpuFamily (INTEL_SKYLAKE | INTEL_XEON; omit to use the location default), availabilityZone (AUTO | ZONE_1 | ZONE_2), plus nics (LAN wiring, firewall, flow logs) and volumes (boot image, size in GB, SSH keys, user-data).

Note: every scale-out provisions fresh volumes. If deleteVolumes is false they are never reclaimed automatically.`,
		Example: fmt.Sprintf(`# Print an editable template of the full properties object, then create from a file
ionosctl vm-autoscaling group create %s > group.json
ionosctl vm-autoscaling group create --json-properties group.json

# Create directly from an inline properties object (CPU-based, 1-10 replicas, add/remove 1 VM per action)
ionosctl vm-autoscaling group create --json-properties '{"properties":{"datacenter":{"id":"<datacenter-id>"},"name":"web-tier","minReplicaCount":1,"maxReplicaCount":10,"policy":{"metric":"INSTANCE_CPU_UTILIZATION_AVERAGE","range":"2m","unit":"PER_MINUTE","scaleInThreshold":33,"scaleOutThreshold":77,"scaleInAction":{"amount":1,"amountType":"ABSOLUTE","cooldownPeriod":"5m","terminationPolicy":"OLDEST_SERVER_FIRST","deleteVolumes":true},"scaleOutAction":{"amount":1,"amountType":"ABSOLUTE","cooldownPeriod":"5m"}},"replicaConfiguration":{"availabilityZone":"AUTO","cores":2,"cpuFamily":"INTEL_SKYLAKE","ram":2048,"nics":[{"lan":1,"name":"nic1","dhcp":true}],"volumes":[{"imageAlias":"ubuntu:latest","name":"boot","size":30,"type":"SSD","imagePassword":"<password>"}]}}}'`,
			core.FlagsUsage(constants.FlagJsonPropertiesExample)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagJsonProperties},
				[]string{constants.FlagJsonPropertiesExample},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			group, _, err := client.Must().VMAscClient.GroupsPost(context.Background()).
				GroupPost(jsonStruct).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(group)
		},
	})

	return cmd
}

var exampleJson = "{\n  \"properties\": {\n    \"datacenter\": {\n      \"id\": \"6e928bd0-3a8e-4821-a20a-54984b0c2d21\"\n    },\n    \"maxReplicaCount\": 10,\n    \"minReplicaCount\": 1,\n    \"name\": \"AutoScaling-Group1\",\n    \"policy\": {\n      \"metric\": \"INSTANCE_CPU_UTILIZATION_AVERAGE\",\n      \"range\": \"2m\",\n      \"scaleInAction\": {\n        \"amount\": 1,\n        \"amountType\": \"ABSOLUTE\",\n        \"cooldownPeriod\": \"5m\",\n        \"terminationPolicy\": \"OLDEST_SERVER_FIRST\",\n        \"deleteVolumes\": true\n      },\n      \"scaleInThreshold\": 33,\n      \"scaleOutAction\": {\n        \"amount\": 1,\n        \"amountType\": \"ABSOLUTE\",\n        \"cooldownPeriod\": \"5m\"\n      },\n      \"scaleOutThreshold\": 77,\n      \"unit\": \"PER_MINUTE\"\n    },\n    \"replicaConfiguration\": {\n      \"availabilityZone\": \"AUTO\",\n      \"cores\": 2,\n      \"cpuFamily\": \"INTEL_SKYLAKE\",\n      \"nics\": [\n        {\n          \"lan\": 1,\n          \"name\": \"LAN-NIC-1\",\n          \"dhcp\": true,\n          \"firewallActive\": true,\n          \"firewallType\": \"INGRESS\",\n          \"flowLogs\": [\n            {\n              \"name\": \"flow-log\",\n              \"action\": \"ACCEPTED\",\n              \"direction\": \"EGRESS\",\n              \"bucket\": \"bucketName/key\"\n            }\n          ],\n          \"firewallRules\": [\n            {\n              \"name\": \"My-resource\",\n              \"protocol\": \"TCP\",\n              \"sourceMac\": \"00:0a:95:9d:68:16\",\n              \"sourceIp\": \"22.231.113.64\",\n              \"targetIp\": \"22.231.113.64\",\n              \"icmpCode\": 2,\n              \"icmpType\": 8,\n              \"portRangeStart\": 8,\n              \"portRangeEnd\": 8,\n              \"type\": \"INGRESS\"\n            }\n          ],\n          \"targetGroup\": {\n            \"targetGroupId\": \"6e928bd0-3a8e-4821-a20a-54984b0c2d21\",\n            \"port\": 8080,\n            \"weight\": 15\n          }\n        }\n      ],\n      \"ram\": 2048,\n      \"volumes\": [\n        {\n          \"image\": \"6e928bd0-3a8e-4821-a20a-54984b0c2d21\",\n          \"imageAlias\": \"ubuntu:latest\",\n          \"name\": \"Volume-1\",\n          \"size\": 30,\n          \"sshKeys\": [\n            \"ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSU\\nGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3\\nPbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XA\\nt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/En\\nmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbx\\nNrRFi9wrf+M7Q== user@domain.local\"\n          ],\n          \"type\": \"SSD\",\n          \"userData\": \"ZWNobyAiSGVsbG8sIFdvcmxkIgo=\",\n          \"bus\": \"VIRTIO\",\n          \"backupunitId\": \"25f67991-0f51-4efc-a8ad-ef1fb31a481c\",\n          \"bootOrder\": \"AUTO\",\n          \"imagePassword\": \"passw0rd\"\n        }\n      ]\n    }\n  }\n}\n"
