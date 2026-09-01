---
description: "List the scaling actions (SCALE_IN / SCALE_OUT history) of one or all groups"
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

List the scaling actions recorded for a group - the log of every SCALE_IN and SCALE_OUT it performed, each with a status (IN_PROGRESS / SUCCESSFUL / FAILED). Pass --group-id to see the history of one group, or --all to gather actions across every group in your account (this fetches actions group-by-group, so it is slower with many groups).

## Options

```text
  -a, --all               List actions from every VM Auto Scaling group in your account (mutually exclusive with --group-id)
  -u, --api-url string    Override default host URL. Preferred over the config file override 'autoscaling' and env var 'IONOS_API_URL' (default "https://api.ionos.com/autoscaling")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ActionId GroupId]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -i, --group-id string   The ID of the VM Auto Scaling group whose scaling actions to list
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl vm-autoscaling action list --group-id GROUP_ID
ionosctl vm-autoscaling action list --all ALL
```

