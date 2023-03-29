---
description: Create a Load Balancer
---

# LoadbalancerCreate

## Usage

```text
ionosctl loadbalancer create [flags]
```

## Aliases

For `loadbalancer` command:

```text
[lb]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a new Load Balancer within the Virtual Data Center. Load balancers can be used for public or private IP traffic. The name, IP and DHCP for the Load Balancer can be set.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --dhcp                   Indicates if the Load Balancer will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false (default true)
  -n, --name string            Name of the Load Balancer (default "Load Balancer")
  -t, --timeout int            Timeout option for Request for Load Balancer creation [seconds] (default 60)
  -w, --wait-for-request       Wait for Request for Load Balancer creation to be executed
```

## Examples

```text
ionosctl loadbalancer create --datacenter-id DATACENTER_ID --name NAME
```

