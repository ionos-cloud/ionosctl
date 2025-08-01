---
description: "Update an Application Load Balancer"
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
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ApplicationLoadBalancerId Name ListenerLan Ips TargetLan PrivateIps State] (default [ApplicationLoadBalancerId,Name,ListenerLan,Ips,TargetLan,PrivateIps,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --ips strings                         Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.
      --listener-lan int                    ID of the listening (inbound) LAN.
  -n, --name string                         The name of the Application Load Balancer. (default "Application Load Balancer")
      --no-headers                          Don't print table headers when table output is used
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
      --private-ips strings                 Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.
  -q, --quiet                               Quiet output
      --target-lan int                      ID of the balanced private target LAN (outbound).
  -t, --timeout int                         Timeout option for Request for Application Load Balancer update [seconds] (default 300)
  -v, --verbose                             Print step-by-step process when running command
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer update to be executed
```

## Examples

```text
ionosctl applicationloadbalancer update --datacenter-id DATACENTER_ID -i APPLICATIONLOADBALANCER_ID --name NAME
```

