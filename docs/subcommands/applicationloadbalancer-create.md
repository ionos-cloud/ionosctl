---
description: Create a Application Load Balancer
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

Use this command to create a Application Load Balancer in a specified Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ApplicationLoadBalancerId Name ListenerLan Ips TargetLan LbPrivateIps State] (default [ApplicationLoadBalancerId,Name,ListenerLan,Ips,TargetLan,LbPrivateIps,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            Collection of IP addresses of the Application Load Balancer. IP of the Listener Lan must be a customer reserved IP for the public Load Balancer and private IP for the private Load Balancer
      --listener-lan int       Id of the listening LAN (default 2)
  -n, --name string            Name of the Application Load Balancer (default "Unnamed Application Load Balancer")
  -o, --output string          Desired output format [text|json] (default "text")
      --private-ips strings    Collection of private IP addresses with subnet mask of the Application Load Balancer. IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet
  -q, --quiet                  Quiet output
      --target-lan int         Id of the balanced private target LAN (default 1)
  -t, --timeout int            Timeout option for Request for Application Load Balancer creation [seconds] (default 300)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait-for-request       Wait for the Request for Application Load Balancer creation to be executed
```

## Examples

```text
ionosctl applicationloadbalancer create --datacenter-id DATACENTER_ID
```

