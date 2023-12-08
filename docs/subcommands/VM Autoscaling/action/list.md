---
description: "List VM Autoscaling Actions"
---

# VmAutoscalingActionList

## Usage

```text
ionosctl vm-autoscaling action list [flags]
```

## Aliases

For `vm-autoscaling` command:

```text
[vmas vm-as vmasc vm-asc vmautoscaling]
```

For `action` command:

```text
[act]
```

For `list` command:

```text
[l ls]
```

## Description

List VM Autoscaling Actions

## Options

```text
  -a, --all               If set, list all actions of all groups
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ActionId]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32       Controls the detail depth of the response objects (default 1)
  -f, --force             Force command to execute without user input
  -i, --group-id string   ID of the autoscaling group to list servers from
  -h, --help              Print usage
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl vm-autoscaling action list --group-id GROUP_ID
ionosctl vm-autoscaling action list --all ALL
```

