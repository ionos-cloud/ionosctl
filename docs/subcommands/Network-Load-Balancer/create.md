---
description: "Create a Network Load Balancer"
---

# NetworkloadbalancerCreate

## Usage

```text
ionosctl compute networkloadbalancer create [flags]
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

Use this command to create a Network Load Balancer (NLB) in a Virtual Data Center. An NLB is a layer-4 (TCP) load balancer that bridges a listener LAN (where clients connect) and a target LAN (the private network of the backend VMs it balances).

After creation the NLB has no forwarding rules yet - it will not forward any traffic until you add at least one rule with `nlb rule create` and attach targets with `nlb rule target add`.

Networking flags:
  --listener-lan: LAN clients connect through (inbound). Defaults to 2.
  --target-lan:   private LAN of the balanced backend VMs (outbound). Defaults to 1.
  --ips:          the NLB's own IP addresses on the listener LAN. For a PUBLIC NLB these must be customer-reserved public IPs; for a PRIVATE NLB they are private IPs. These are also the IPs a forwarding rule listens on.
  --private-ips:  the NLB's private IPs (with subnet mask, e.g. 10.0.0.5/24) on the target LAN, used to reach the backends. If omitted, the system assigns one with a /24 mask.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NetworkLoadBalancerId Name ListenerLan Ips TargetLan LbPrivateIps State DatacenterId]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            IP addresses of the NLB on the listener LAN, also used as forwarding-rule listener IPs. Must be customer-reserved public IPs for a public NLB, or private IPs for a private NLB
      --limit int              Maximum number of items to return per request (default 50)
      --listener-lan int       ID of the listener LAN (inbound) where clients connect to the NLB (default 2)
  -n, --name string            Name of the Network Load Balancer (default "Network Load Balancer")
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --private-ips strings    Private IP addresses (with subnet mask, e.g. 10.0.0.5/24) the NLB uses on the target LAN to reach backends. If omitted, an IP with a /24 mask is generated
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --target-lan int         ID of the private target LAN (outbound) hosting the backend VMs the NLB balances (default 1)
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create an NLB with the default listener/target LANs
ionosctl compute networkloadbalancer create --datacenter-id DATACENTER_ID --name "web-nlb"

# Create a public NLB, pinning the LANs and giving it a reserved public listener IP
ionosctl compute networkloadbalancer create --datacenter-id DATACENTER_ID --name "prod-nlb" --listener-lan 2 --target-lan 1 --ips 203.0.113.10 --private-ips 10.0.0.5/24
```

