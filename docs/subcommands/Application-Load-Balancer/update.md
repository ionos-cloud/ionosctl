---
description: "Update an Application Load Balancer"
---

# ApplicationloadbalancerUpdate

## Usage

```text
ionosctl compute applicationloadbalancer update [flags]
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

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Application Load Balancer Id

## Options

```text
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ApplicationLoadBalancerId Name ListenerLan Ips TargetLan PrivateIps State DatacenterId]
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int                           Level of detail for response objects (default 1)
  -F, --filters strings                     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --ips strings                         Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.
      --limit int                           Maximum number of items to return per request (default 50)
      --listener-lan int                    ID of the listening (inbound) LAN.
  -n, --name string                         The name of the Application Load Balancer. (default "Application Load Balancer")
      --no-headers                          Don't print table headers when table output is used
      --offset int                          Number of items to skip before starting to collect the results
      --order-by string                     Property to order the results by
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
      --private-ips strings                 Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.
      --query string                        JMESPath query string to filter the output
  -q, --quiet                               Quiet output
      --target-lan int                      ID of the balanced private target LAN (outbound).
  -t, --timeout int                         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                                Wait for the resource to reach AVAILABLE state after the command completes
```

## Examples

```text
ionosctl compute applicationloadbalancer update --datacenter-id DATACENTER_ID -i APPLICATIONLOADBALANCER_ID --name NAME
```

