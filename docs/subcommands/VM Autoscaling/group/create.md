---
description: "Create a VM-Autoscaling group"
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

Create a VM-Autoscaling group

## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [GroupId DatacenterId Name MinReplicas Replicas MaxReplicas Location State Metric Range ScaleInActionAmount ScaleInActionAmountType ScaleInActionCooldownPeriod ScaleInActionTerminationPolicy ScaleInActionDeleteVolumes ScaleInThreshold ScaleOutActionAmount ScaleOutActionAmountType ScaleOutActionCooldownPeriod ScaleOutThreshold Unit AvailabilityZone Cores CPUFamily RAM]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --json-properties string    Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example   If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
  -t, --timeout int               Timeout in seconds for polling the request (default 60)
  -v, --verbose                   Print step-by-step process when running command
  -w, --wait                      Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl vm-autoscaling group create --json-properties JSON_PROPERTIES 
```

