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
  -u, --api-url string     Override default host URL. Preferred over the config file override 'autoscaling' and env var 'IONOS_API_URL' (default "https://api.ionos.com/autoscaling")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [ActionId GroupId]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32        Controls the detail depth of the response objects (default 1)
  -f, --force              Force command to execute without user input
      --group-id string    ID of the autoscaling group that the action is a part of
  -h, --help               Print usage
      --limit int          Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers         Don't print table headers when table output is used
      --offset int         Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -v, --verbose count      Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl vm-autoscaling action get --group-id GROUP_ID --action-id ACTION_ID 
```

