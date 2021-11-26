---
description: List Network Load Balancer FlowLogs
---

# NetworkloadbalancerFlowlogList

## Usage

```text
ionosctl networkloadbalancer flowlog list [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `flowlog` command:

```text
[f fl]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Network Load Balancer FlowLogs from a specified Network Load Balancer.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name action direction bucket]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* Data Center Id
* Network Load Balancer Id

## Options

```text
  -u, --api-url string                  Override default host url (default "https://api.ionos.com")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string            The unique Data Center Id (required)
  -F, --filters strings                 Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
  -M, --max-results int                 The maximum number of elements to return
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --order-by string                 Limits results to those containing a matching value for a specific property
  -o, --output string                   Desired output format [text|json] (default "text")
  -q, --quiet                           Quiet output
  -v, --verbose                         Print step-by-step process when running command
```

## Examples

```text
ionosctl networkloadbalancer flowlog list --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID
```

