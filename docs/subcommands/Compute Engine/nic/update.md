---
description: "Update a NIC"
---

# NicUpdate

## Usage

```text
ionosctl compute nic update [flags]
```

## Aliases

For `nic` command:

```text
[n]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update the configuration of an existing NIC, identified within its Data Center (--datacenter-id) and Server (--server-id) by --nic-id. Only the flags you set are changed; everything else is left as-is.

Common updates:
* Move the NIC to a different network by changing --lan-id.
* Add reserved public IPs or assign private IPs (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) with --ips.
* Toggle DHCP with --dhcp.
* Enable/disable the per-NIC firewall with --firewall-active and set its direction with --firewall-type (INGRESS/EGRESS/BIDIRECTIONAL). When enabled, incoming traffic is filtered by the NIC's firewall rules; when disabled, all traffic reaches the NIC directly and the rules are ignored.

Restriction: the primary address of a NIC connected to a Load Balancer can only be changed by changing the IP of the Load Balancer.

Use `--wait` (`-w`) to wait for the NIC to reach AVAILABLE state before the command returns.

Required values to run command:

* Data Center Id
* NIC Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Name Dhcp LanId Ips IPv6Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac DHCPv6 IPv6CidrBlock]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
      --dhcp                   Whether the NIC reserves an IP automatically from the LAN's DHCP server (true) or not (false). Set --dhcp=false when managing addresses via --ips (default true)
      --dhcpv6                 Whether the NIC reserves an IPv6 address automatically via DHCPv6. Only applies when the target LAN has IPv6 enabled. Set --dhcpv6=false to disable (default true)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
      --firewall-active        Enable the per-NIC firewall. When enabled, an empty ruleset blocks all incoming traffic and only explicitly configured firewall rules are allowed through; when disabled, all traffic reaches the NIC and rules are ignored
      --firewall-type string   Direction of traffic the NIC's firewall rules apply to. INGRESS = inbound only, EGRESS = outbound only, BIDIRECTIONAL = both. Only meaningful when --firewall-active=true (default "INGRESS")
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            IPs to assign to the NIC. Explicitly assigned public IPs must come from a reserved IP block; private-range IPs (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) may be used on private LANs
      --ipv6-cidr string       The /80 IPv6 Cidr as defined in RFC 4291. It needs to be within the LAN IPv6 Cidr Block range. (default "disable")
      --ipv6-ips strings       One or more IPv6 IPs to assign to the NIC. Each must fall within the NIC's IPv6 CIDR block, and the target LAN must have IPv6 enabled
      --lan-id int             Move the NIC to this LAN ID, changing which network the server reaches through this NIC. If the LAN ID does not exist in the Data Center, it is created implicitly (default 1)
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            A human-friendly name for the NIC (shown in the DCD and listings). Does not affect networking
  -i, --nic-id string          The unique NIC Id (required)
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
# Move a NIC to a different LAN
ionosctl compute nic update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --lan-id 2 --wait

# Rename a NIC, add a reserved public IP, and enable a bidirectional firewall
ionosctl compute nic update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --name web-nic --ips 203.0.113.10 --firewall-active=true --firewall-type BIDIRECTIONAL --wait
```

