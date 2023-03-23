---
description: Get a Load Balancer
---

# LoadbalancerGet

## Usage

```text
ionosctl loadbalancer get [flags]
```

## Aliases

For `loadbalancer` command:

```text
[lb]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information about a Load Balancer instance.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int32              Controls the detail depth of the response objects. Max depth is 10.
  -i, --loadbalancer-id string   The unique Load Balancer Id (required)
      --no-headers               When using text output, don't print headers
```

## Examples

```text
ionosctl loadbalancer get --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID
```

