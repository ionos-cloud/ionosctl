---
description: "List VM Autoscaling Groups. Use a greater '--depth' to see current replica count"
---

# VmAutoscalingGroupList

## Usage

```text
ionosctl vm-autoscaling group list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

List VM Autoscaling Groups. Use a greater '--depth' to see current replica count

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [GroupId DatacenterId Name MinReplicas Replicas MaxReplicas Location State Metric Range ScaleInActionAmount ScaleInActionAmountType ScaleInActionCooldownPeriod ScaleInActionTerminationPolicy ScaleInActionDeleteVolumes ScaleInThreshold ScaleOutActionAmount ScaleOutActionAmountType ScaleOutActionCooldownPeriod ScaleOutThreshold Unit AvailabilityZone Cores CPUFamily RAM]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32      Controls the detail depth of the response objects (default 1)
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl vm-autoscaling group list
```

