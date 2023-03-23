---
description: Delete a Load Balancer
---

# LoadbalancerDelete

## Usage

```text
ionosctl loadbalancer delete [flags]
```

## Aliases

For `loadbalancer` command:

```text
[lb]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete the specified Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
  -a, --all                      Delete all the LoadBlancers from a virtual Datacenter.
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int32              Controls the detail depth of the response objects. Max depth is 10.
  -i, --loadbalancer-id string   The unique Load Balancer Id (required)
  -t, --timeout int              Timeout option for Request for Load Balancer deletion [seconds] (default 60)
  -w, --wait-for-request         Wait for Request for Load Balancer deletion to be executed
```

## Examples

```text
ionosctl loadbalancer delete --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --force --wait-for-request
```

