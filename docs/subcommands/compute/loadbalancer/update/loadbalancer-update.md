---
description: Update a Load Balancer
---

# LoadbalancerUpdate

## Usage

```text
ionosctl loadbalancer update [flags]
```

## Aliases

For `loadbalancer` command:

```text
[lb]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update the configuration of a specified Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Load Balancer Id

## Options

```text
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int32              Controls the detail depth of the response objects. Max depth is 10.
      --dhcp                     Indicates if the Load Balancer will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false (default true)
      --ip ip                    The IP of the Load Balancer
  -i, --loadbalancer-id string   The unique Load Balancer Id (required)
  -n, --name string              Name of the Load Balancer
  -t, --timeout int              Timeout option for Request for Load Balancer update [seconds] (default 60)
  -w, --wait-for-request         Wait for Request for Load Balancer update to be executed
```

## Examples

```text
ionosctl loadbalancer update --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --dhcp=false --wait-for-request
```

