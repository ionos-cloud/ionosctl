---
description: Create a Network Load Balancer
---

# NetworkloadbalancerCreate

## Usage

```text
ionosctl networkloadbalancer create [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a Network Load Balancer in a specified Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --ips strings            Collection of IP addresses of the Network Load Balancer
      --listener-lan int       Id of the listening LAN (default 2)
  -n, --name string            Name of the Network Load Balancer (default "Network Load Balancer")
      --private-ips strings    Collection of private IP addresses with subnet mask of the Network Load Balancer
      --target-lan int         Id of the balanced private target LAN (default 1)
  -t, --timeout int            Timeout option for Request for Network Load Balancer creation [seconds] (default 300)
  -w, --wait-for-request       Wait for the Request for Network Load Balancer creation to be executed
```

## Examples

```text
ionosctl networkloadbalancer create --datacenter-id DATACENTER_ID
```

