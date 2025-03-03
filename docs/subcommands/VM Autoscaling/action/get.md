---
description: "Get a VM Autoscaling Action"
---

# VmAutoscalingActionGet

## Usage

```text
ionosctl vm-autoscaling action get [flags]
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

For `get` command:

```text
[g]
```

## Description

Get a VM Autoscaling Action

## Options

```text
  -i, --action-id string   ID of the autoscaling action
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [ActionId GroupId]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32        Controls the detail depth of the response objects (default 1)
  -f, --force              Force command to execute without user input
      --group-id string    ID of the autoscaling group that the action is a part of
  -h, --help               Print usage
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -v, --verbose count      Print step-by-step process when running command
```

## Examples

```text
ionosctl vm-autoscaling action get --group-id GROUP_ID --action-id ACTION_ID 
```

