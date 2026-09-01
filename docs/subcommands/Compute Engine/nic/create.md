---
description: "Create a NIC"
---

# NicCreate

## Usage

```text
ionosctl compute nic create [flags]
```

## Aliases

For `nic` command:

```text
[n]
```

For `create` command:

```text
[c]
```

## Description

Use this command to add a new NIC (Network Interface Card) to a server. The NIC attaches the server (--server-id) to a LAN (--lan-id) inside the given Data Center (--datacenter-id). If the target LAN does not exist yet, it is created implicitly when the NIC is created.

Addressing options:
* DHCP (default): with --dhcp=true the NIC reserves an IP automatically from the LAN's DHCP server. This is the usual choice.
* Static/reserved IPs: pass --ips to assign specific addresses. Explicitly assigned public IPs must come from reserved IP blocks; private-range addresses (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) may be assigned on private LANs. Leaving --ips empty lets IONOS pick an address automatically.

Firewall: --firewall-active toggles the per-NIC firewall (off by default). When active, all incoming traffic is blocked except what explicit firewall rules allow; when inactive, rules are ignored and traffic reaches the NIC directly. --firewall-type selects which traffic direction those rules govern (INGRESS/EGRESS/BIDIRECTIONAL).

IPv6: --ipv6-cidr-block, --ipv6-ips and --dhcpv6 only apply when the target LAN has IPv6 enabled; --ipv6-ips must fall within the NIC's IPv6 CIDR block.

Use `--wait` (`-w`) to wait for the NIC to reach AVAILABLE state before the command returns.

Required values to run a command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Name Dhcp LanId Ips IPv6Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac DHCPv6 IPv6CidrBlock]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
      --dhcp                   Whether the NIC reserves an IP automatically from the LAN's DHCP server. Set --dhcp=false to disable DHCP (typically when assigning static --ips). Default: true (default true)
      --dhcpv6                 Whether the NIC reserves an IPv6 address automatically via DHCPv6. Only applies when the target LAN has IPv6 enabled. Set --dhcpv6=false to disable. Default: true (default true)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
      --firewall-active        Enable the per-NIC firewall. When enabled, an empty ruleset blocks all incoming traffic and only explicitly configured firewall rules are allowed through; when disabled, all traffic reaches the NIC and rules are ignored. Default: false
      --firewall-type string   Direction of traffic the NIC's firewall rules apply to. INGRESS = inbound only, EGRESS = outbound only, BIDIRECTIONAL = both. Only meaningful when --firewall-active=true. Default: INGRESS (default "INGRESS")
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            One or more IPs to assign to the NIC. Explicitly assigned public IPs must come from a reserved IP block; private-range IPs (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) may be used on private LANs. Leave empty to let IONOS assign an address automatically (see --dhcp)
      --ipv6-cidr string       The /80 IPv6 Cidr as defined in RFC 4291. It needs to be within the LAN IPv6 Cidr Block range. (default "disable")
      --ipv6-ips strings       One or more IPv6 IPs to assign to the NIC. Each must fall within the NIC's IPv6 CIDR block (--ipv6-cidr-block), and the target LAN must have IPv6 enabled
      --lan-id int             The ID of the LAN this NIC attaches to, determining which network the server reaches through this NIC. If the LAN ID does not exist in the Data Center, it is created implicitly. Default: 1 (default 1)
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            A human-friendly name for the NIC (shown in the DCD and listings). Does not affect networking (default "Internet Access")
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a DHCP NIC on LAN 1 (defaults: --dhcp=true, --lan-id=1, firewall off)
ionosctl compute nic create --datacenter-id DATACENTER_ID --server-id SERVER_ID --name mynic

# Create a NIC on a specific LAN with static reserved IPs and an active ingress firewall
ionosctl compute nic create --datacenter-id DATACENTER_ID --server-id SERVER_ID --name web-nic --lan-id 2 --dhcp=false --ips 203.0.113.10 --firewall-active=true --firewall-type INGRESS --wait
```

