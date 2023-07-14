---
description: "Create an Application Load Balancer"
---

# ApplicationloadbalancerCreate

## Usage

```text
ionosctl applicationloadbalancer create [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create an Application Load Balancer in a specified Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ApplicationLoadBalancerId Name ListenerLan Ips TargetLan PrivateIps State] (default [ApplicationLoadBalancerId,Name,ListenerLan,Ips,TargetLan,PrivateIps,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.
      --listener-lan int       ID of the listening (inbound) LAN. (default 2)
  -n, --name string            The name of the Application Load Balancer. (default "Unnamed Application Load Balancer")
  -o, --output string          Desired output format [text|json] (default "text")
      --private-ips strings    Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.
  -q, --quiet                  Quiet output
      --target-lan int         ID of the balanced private target LAN (outbound). (default 1)
  -t, --timeout int            Timeout option for Request for Application Load Balancer creation [seconds] (default 10000)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait-for-request       Wait for the Request for Application Load Balancer creation to be executed
```

## Examples

```text
ionosctl applicationloadbalancer create --datacenter-id DATACENTER_ID
```

