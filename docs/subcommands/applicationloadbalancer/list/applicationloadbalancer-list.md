---
description: List Application Load Balancers
---

# ApplicationloadbalancerList

## Usage

```text
ionosctl applicationloadbalancer list [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Application Load Balancers from a specified Virtual Data Center.

Required values to run command:

* Data Center Id

## Options

```text
  -a, --all                    List all resources without the need of specifying parent ID name.
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32      The maximum number of elements to return
      --order-by string        Limits results to those containing a matching value for a specific property
```

## Examples

```text
ionosctl applicationloadbalancer list --datacenter-id DATACENTER_ID
```

