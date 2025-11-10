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
  -u, --api-url string            Override default host URL. Preferred over the config file override 'autoscaling' and env var 'IONOS_API_URL' (default "https://api.ionos.com/autoscaling")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [GroupId DatacenterId Name MinReplicas Replicas MaxReplicas Location State Metric Range ScaleInActionAmount ScaleInActionAmountType ScaleInActionCooldownPeriod ScaleInActionTerminationPolicy ScaleInActionDeleteVolumes ScaleInThreshold ScaleOutActionAmount ScaleOutActionAmountType ScaleOutActionCooldownPeriod ScaleOutThreshold Unit AvailabilityZone Cores CPUFamily RAM]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --json-properties string    Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example   If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
      --limit int                 Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers                Don't print table headers when table output is used
      --offset int                Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl vm-autoscaling group create --json-properties JSON_PROPERTIES 
```

