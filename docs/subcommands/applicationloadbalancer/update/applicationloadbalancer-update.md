---
description: Update an Application Load Balancer
---

# ApplicationloadbalancerUpdate

## Usage

```text
ionosctl applicationloadbalancer update [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified Application Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` or `-w` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id

## Options

```text
  -i, --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
      --ips strings                         Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.
      --listener-lan int                    ID of the listening (inbound) LAN.
  -n, --name string                         The name of the Application Load Balancer. (default "Application Load Balancer")
      --private-ips strings                 Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.
      --target-lan int                      ID of the balanced private target LAN (outbound).
  -t, --timeout int                         Timeout option for Request for Application Load Balancer update [seconds] (default 300)
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer update to be executed
```

## Examples

```text
ionosctl applicationloadbalancer update --datacenter-id DATACENTER_ID -i APPLICATIONLOADBALANCER_ID --name NAME
```

