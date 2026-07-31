---
description: "Create a VM Auto Scaling group from a JSON properties object"
---

# VmAutoscalingGroupCreate

## Usage

```text
ionosctl vm-autoscaling group create [flags]
```

## Aliases

For `vm-autoscaling` command:

```text
[vmas vm-as vmasc vm-asc vmautoscaling]
```

For `group` command:

```text
[g groups]
```

For `create` command:

```text
[c]
```

## Description

Create a new VM Auto Scaling group. The group's full configuration is supplied as a JSON 'properties' object (--json-properties, or a path to a file), because it nests three sub-objects that would be unwieldy as individual flags. Run this command with --json-properties-example to print a ready-to-edit template.

The properties object has these parts:

  * datacenter.id (required) - the datacenter the group and all its replicas live in.
  * name - a friendly label for the group.
  * minReplicaCount / maxReplicaCount (0-100) - the floor and ceiling for the number of running replicas. The autoscaler stays within these bounds.
  * policy - when and how to scale:
      - metric: INSTANCE_CPU_UTILIZATION_AVERAGE | INSTANCE_NETWORK_IN_BYTES | INSTANCE_NETWORK_IN_PACKETS | INSTANCE_NETWORK_OUT_BYTES | INSTANCE_NETWORK_OUT_PACKETS
      - range: ISO 8601 duration for the metric aggregation window (e.g. P1D = 1 day, PT30M = 30 min). Must be >= 2 minutes; API default 120s.
      - unit: PER_SECOND | PER_MINUTE | PER_HOUR | TOTAL - how the metric is normalized. For network metrics this is the rate unit; if unit=TOTAL, scaleOutThreshold must be >= 40.
      - scaleOutThreshold: when the metric goes ABOVE this value, scale out. scaleInThreshold: when it drops BELOW this value, scale in. The two thresholds must keep a minimum gap (metric-dependent) so the group does not scale in and out at the same time. scaleInThreshold < scaleOutThreshold.
      - scaleOutAction / scaleInAction: amount = how many replicas (or what percentage, when amountType=PERCENTAGE) to add/remove; amountType = ABSOLUTE | PERCENTAGE; cooldownPeriod = how long to wait before the next action (min 2m, max 24h, default 5m).
      - scaleInAction.terminationPolicy: OLDEST_SERVER_FIRST | NEWEST_SERVER_FIRST | RANDOM - which replica to remove first. scaleInAction.deleteVolumes: if true, a removed replica's volumes are deleted too. Leave true unless you need to keep the data; orphaned volumes count against your contract limits and can eventually block further scale-outs.
  * replicaConfiguration - the template every new replica is cloned from: cores, ram (in MB), cpuFamily (INTEL_SKYLAKE | INTEL_XEON; omit to use the location default), availabilityZone (AUTO | ZONE_1 | ZONE_2), plus nics (LAN wiring, firewall, flow logs) and volumes (boot image, size in GB, SSH keys, user-data).

Note: every scale-out provisions fresh volumes. If deleteVolumes is false they are never reclaimed automatically.

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'autoscaling' and env var 'IONOS_API_URL' (default "https://api.ionos.com/autoscaling")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [GroupId DatacenterId Name MinReplicas Replicas MaxReplicas Location State Metric Range ScaleInActionAmount ScaleInActionAmountType ScaleInActionCooldownPeriod ScaleInActionTerminationPolicy ScaleInActionDeleteVolumes ScaleInThreshold ScaleOutActionAmount ScaleOutActionAmountType ScaleOutActionCooldownPeriod ScaleOutThreshold Unit AvailabilityZone Cores CPUFamily RAM]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --json-properties string    Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example   If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
      --limit int                 Maximum number of items to return per request (default 50)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Print an editable template of the full properties object, then create from a file
ionosctl vm-autoscaling group create --json-properties-example JSON_PROPERTIES_EXAMPLE  > group.json
ionosctl vm-autoscaling group create --json-properties group.json

# Create directly from an inline properties object (CPU-based, 1-10 replicas, add/remove 1 VM per action)
ionosctl vm-autoscaling group create --json-properties '{"properties":{"datacenter":{"id":"<datacenter-id>"},"name":"web-tier","minReplicaCount":1,"maxReplicaCount":10,"policy":{"metric":"INSTANCE_CPU_UTILIZATION_AVERAGE","range":"PT2M","unit":"PER_MINUTE","scaleInThreshold":33,"scaleOutThreshold":77,"scaleInAction":{"amount":1,"amountType":"ABSOLUTE","cooldownPeriod":"5m","terminationPolicy":"OLDEST_SERVER_FIRST","deleteVolumes":true},"scaleOutAction":{"amount":1,"amountType":"ABSOLUTE","cooldownPeriod":"5m"}},"replicaConfiguration":{"availabilityZone":"AUTO","cores":2,"cpuFamily":"INTEL_SKYLAKE","ram":2048,"nics":[{"lan":1,"name":"nic1","dhcp":true}],"volumes":[{"imageAlias":"ubuntu:latest","name":"boot","size":30,"type":"SSD","imagePassword":"<password>"}]}}}'
```

