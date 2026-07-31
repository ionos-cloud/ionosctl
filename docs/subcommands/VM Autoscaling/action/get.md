---
description: "Get a single VM Auto Scaling action"
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

Show one scaling action of a group by its ID. Because actions are scoped to a group, both --group-id (which group the action belongs to) and --action-id (the action itself) are required. The result includes the action type (SCALE_IN / SCALE_OUT) and its status (IN_PROGRESS / SUCCESSFUL / FAILED) - useful for confirming whether a specific scaling event completed.

## Options

```text
  -i, --action-id string   The ID of the scaling action to show (must belong to --group-id)
  -u, --api-url string     Override default host URL. Preferred over the config file override 'autoscaling' and env var 'IONOS_API_URL' (default "https://api.ionos.com/autoscaling")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [ActionId GroupId]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int          Level of detail for response objects (default 1)
  -F, --filters strings    Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force              Force command to execute without user input
      --group-id string    The ID of the VM Auto Scaling group the action belongs to
  -h, --help               Print usage
      --limit int          Maximum number of items to return per request (default 50)
      --no-headers         Don't print table headers when table output is used
      --offset int         Number of items to skip before starting to collect the results
      --order-by string    Property to order the results by
  -o, --output string      Desired output format [text|json|api-json] (default "text")
      --query string       JMESPath query string to filter the output
  -q, --quiet              Quiet output
  -t, --timeout int        Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count      Increase verbosity level [-v, -vv, -vvv]
  -w, --wait               Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl vm-autoscaling action get --group-id GROUP_ID --action-id ACTION_ID 
```

