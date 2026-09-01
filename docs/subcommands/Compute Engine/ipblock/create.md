---
description: "Reserve a block of static public IPv4 addresses"
---

# IpblockCreate

## Usage

```text
ionosctl compute ipblock create [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `create` command:

```text
[c]
```

## Description

Reserve an IpBlock: a set of `--size` static, public IPv4 addresses held in one `--location` (region). The reserved IPs can then be assigned to NICs, NAT gateways, load balancers or IP-failover groups in any Virtual Data Center in that SAME location.

Both --location and --size are fixed at reservation time and cannot be changed later - a block cannot be moved to another region or resized. Reserved IPs are billed while held and persist across server power-offs (unlike dynamic DHCP addresses), so reserve only what you need.

Reservation is asynchronous. Use `--wait` (`-w`) to block until the IpBlock reaches AVAILABLE state and its IPs are ready to assign.

## Options

```text
  -u, --api-url string    Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [IpBlockId Name Location Size Ips State]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -l, --location string   Region the IPs are reserved in, as <country>/<city> (e.g. de/txl Berlin, de/fra Frankfurt, us/las Las Vegas). The block can only serve resources in this region and cannot be moved afterwards. Location de/fra/2 is currently unavailable (default "de/txl")
  -n, --name string       A friendly label for the block (shown in listings; not the IP addresses themselves). If omitted, the API assigns a name automatically
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
      --size int          How many static public IPv4 addresses to reserve in this block. Fixed at creation - a block cannot be resized later. Each reserved IP is billed while held (default 2)
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Reserve a block of 1 IP in Berlin (de/txl)
ionosctl compute ipblock create --name web-vip --location de/txl --size 1 --wait

# Reserve 4 IPs in Frankfurt for a load balancer / failover pool, then list the assigned addresses
ionosctl compute ipblock create --name lb-pool --location de/fra --size 4 --wait
```

