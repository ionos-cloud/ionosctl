---
description: "List Load Balancers"
---

# LoadbalancerList

## Usage

```text
ionosctl loadbalancer list [flags]
```

## Aliases

For `loadbalancer` command:

```text
[lb]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of Load Balancers within a Virtual Data Center on your account.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [dhcp ip name]
* filter by metadata: [createdBy createdByUserId createdDate etag lastModifiedBy lastModifiedByUserId lastModifiedDate state]

Required values to run command:

* Data Center Id

## Options

```text
  -a, --all                    List all resources without the need of specifying parent ID name.
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LoadBalancerId Name Dhcp State Ip DatacenterId] (default [LoadBalancerId,Name,Dhcp,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -h, --help                   Print usage
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             Don't print table headers when table output is used
      --order-by string        Limits results to those containing a matching value for a specific property
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl loadbalancer list --datacenter-id DATACENTER_ID
```

