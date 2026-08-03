---
description: "Create an Application Load Balancer"
---

# ApplicationloadbalancerCreate

## Usage

```text
ionosctl compute applicationloadbalancer create [flags]
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

Use this command to create a layer-7 Application Load Balancer in a specified Virtual Data Center.

An ALB bridges two LANs in the same data center:
  * The listener LAN (--listener-lan) is where clients reach the balancer. For a public ALB this LAN is internet-facing and --ips holds customer-reserved public IPs; for a private ALB it is an internal LAN and --ips holds private IPs.
  * The target LAN (--target-lan) is the private LAN where your backend servers (grouped into target groups) live. The ALB uses --private-ips as its own addresses on this LAN to reach the backends.

After creating the ALB you attach forwarding rules (`alb rule create`), then HTTP rules within them (`alb rule httprule add`) to route traffic to target groups.

Use `--wait` (`-w`) to block until the ALB reaches AVAILABLE state before its rules are configured.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ApplicationLoadBalancerId Name ListenerLan Ips TargetLan PrivateIps State DatacenterId]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            The IP addresses clients use to reach the balancer on the listener LAN. These are customer-reserved public IPs for a public ALB, or private IPs for a private ALB. Provide one or more, e.g. --ips 192.0.2.10,192.0.2.11
      --limit int              Maximum number of items to return per request (default 50)
      --listener-lan int       Numeric ID of the LAN clients connect to (the inbound/listener LAN). For a public ALB this is an internet-facing LAN; for a private ALB it is an internal LAN. Defaults to 2. (default 2)
  -n, --name string            The name of the Application Load Balancer. (default "Unnamed Application Load Balancer")
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --private-ips strings    The balancer's own private IP addresses (with subnet mask) on the target LAN, used to reach the backends. Each value must include a valid subnet mask, e.g. --private-ips 10.0.1.5/24. If omitted, the system auto-generates an IP with a /24 subnet.
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --target-lan int         Numeric ID of the private LAN where the balanced backend servers live (the outbound/target LAN). Defaults to 1. (default 1)
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a public ALB listening on LAN 2, balancing to backends on LAN 1
ionosctl compute applicationloadbalancer create --datacenter-id DATACENTER_ID --name "web-alb" --listener-lan 2 --target-lan 1 --ips 192.0.2.10

# Create an ALB and wait for it to become AVAILABLE, letting the system auto-assign a /24 private IP on the target LAN
ionosctl compute applicationloadbalancer create --datacenter-id DATACENTER_ID --name "web-alb" --listener-lan 2 --target-lan 1 --ips 192.0.2.10 --private-ips 10.0.1.5/24 --wait
```

