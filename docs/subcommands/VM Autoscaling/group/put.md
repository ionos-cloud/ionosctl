---
description: "Replace the configuration of an existing VM Auto Scaling group"
---

# VmAutoscalingGroupPut

## Usage

```text
ionosctl vm-autoscaling group put [flags]
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

For `put` command:

```text
[p update]
```

## Description

Replace the properties of an existing VM Auto Scaling group (HTTP PUT). This is a full replacement, not a patch: the JSON 'properties' object you pass becomes the group's new configuration in its entirety, so include every field you want to keep. A common workflow is to 'group get' the current group, edit the properties, and pass them back here.

Typical reasons to update a group:
  * Change the replica bounds (minReplicaCount / maxReplicaCount) to widen or narrow how far it can scale. Lowering maxReplicaCount below the current replica count triggers scale-in; raising minReplicaCount above it triggers scale-out.
  * Retune the policy - e.g. adjust scaleInThreshold / scaleOutThreshold, switch the metric or unit, or change the scale action amount and cooldownPeriod.
  * Update the replicaConfiguration template (image, cores, ram, NICs). Note this affects replicas created AFTER the change; existing replicas are not re-provisioned.

See the field-by-field reference under 'group create' for the meaning and valid values of each property.

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'autoscaling' and env var 'IONOS_API_URL' (default "https://api.ionos.com/autoscaling")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [GroupId DatacenterId Name MinReplicas Replicas MaxReplicas Location State Metric Range ScaleInActionAmount ScaleInActionAmountType ScaleInActionCooldownPeriod ScaleInActionTerminationPolicy ScaleInActionDeleteVolumes ScaleInThreshold ScaleOutActionAmount ScaleOutActionAmountType ScaleOutActionCooldownPeriod ScaleOutThreshold Unit AvailabilityZone Cores CPUFamily RAM]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -i, --group-id string           The ID of the VM Auto Scaling group whose configuration is replaced (required)
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
# Widen an existing group's replica bounds and raise its scale-out threshold
ionosctl vm-autoscaling group put --group-id GROUP_ID --json-properties '{"properties":{"datacenter":{"id":"<datacenter-id>"},"name":"web-tier","minReplicaCount":2,"maxReplicaCount":20,"policy":{"metric":"INSTANCE_CPU_UTILIZATION_AVERAGE","range":"2m","unit":"PER_MINUTE","scaleInThreshold":30,"scaleOutThreshold":85,"scaleInAction":{"amount":1,"amountType":"ABSOLUTE","cooldownPeriod":"5m","terminationPolicy":"OLDEST_SERVER_FIRST","deleteVolumes":true},"scaleOutAction":{"amount":2,"amountType":"ABSOLUTE","cooldownPeriod":"10m"}},"replicaConfiguration":{"availabilityZone":"AUTO","cores":2,"cpuFamily":"INTEL_SKYLAKE","ram":2048,"nics":[{"lan":1,"name":"nic1","dhcp":true}],"volumes":[{"imageAlias":"ubuntu:latest","name":"boot","size":30,"type":"SSD","imagePassword":"<password>"}]}}}'
```

