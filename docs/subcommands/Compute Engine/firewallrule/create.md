---
description: "Create a Firewall Rule"
---

# FirewallruleCreate

## Usage

```text
ionosctl compute firewallrule create [flags]
```

## Aliases

For `firewallrule` command:

```text
[f fr firewall]
```

For `create` command:

```text
[c]
```

## Description

Add a new Firewall Rule to the NIC identified by --datacenter-id / --server-id / --nic-id. Every Firewall Rule belongs to exactly one NIC.

A rule WHITELISTS a slice of traffic: while the NIC's firewall is active, traffic is only allowed if a rule matches it (default-deny). --direction (alias --type) selects INGRESS (traffic entering the NIC) or EGRESS (traffic leaving the NIC); it defaults to INGRESS.

--protocol determines which other match flags apply:
  * TCP / UDP  -> --port-range-start and --port-range-end restrict the destination port range. Leave both unset to allow all ports.
  * ICMP       -> --icmp-type and --icmp-code restrict the ICMP message. Leave unset to allow all types/codes.
  * ANY        -> matches every protocol; port and ICMP flags do not apply.
--source-mac, --source-ip and --destination-ip narrow the match further for any protocol; any of them left unset acts as a wildcard (allow all).

NOTE: --protocol is fixed at creation time. It cannot be changed later (the update command has no --protocol flag); to change protocol you must delete the rule and create a new one.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Protocol

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FirewallRuleId Name Protocol PortRangeStart PortRangeEnd Direction IPVersion State SourceMac SourceIP DestinationIP IcmpCode IcmpType]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
      --destination-ip -D      When the NIC has multiple IPs, match only traffic directed to this IP address or CIDR block of the NIC (must match --ip-version). Leave unset to allow any target IP. WARNING: the short-hand flag -D is deprecated
  -d, --direction string       Direction of traffic the rule matches: INGRESS (entering the NIC) or EGRESS (leaving the NIC). Defaults to INGRESS (default "INGRESS")
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --icmp-code int          Only when --protocol ICMP: match this ICMP code (0-254). Leave unset to allow all codes
      --icmp-type int          Only when --protocol ICMP: match this ICMP type (0-254), e.g. 8 = echo request (ping), 0 = echo reply. Leave unset to allow all types
      --ip-version string      The IP version this rule applies to. If --source-ip/--destination-ip are given it must match their version; if omitted it is deduced from those addresses. With no IPs given the rule only allows the selected version (defaults to IPv4). Can be one of: IPv4, IPv6 (default "IPv4")
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            A human-friendly label for the rule. Has no effect on matching; used only to identify the rule in listings (default "Unnamed Rule")
      --nic-id string          The unique NIC Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --port-range-end int     Only when --protocol TCP or UDP: last port of the allowed destination-port range (1-65535, inclusive). Set both --port-range-start and --port-range-end; leave both unset to allow all ports (default 1)
      --port-range-start int   Only when --protocol TCP or UDP: first port of the allowed destination-port range (1-65535, inclusive). Set both --port-range-start and --port-range-end (use the same value for a single port); leave both unset to allow all ports (default 1)
      --protocol string        The IP protocol this rule matches. TCP/UDP also honour --port-range-start/--port-range-end; ICMP also honours --icmp-type/--icmp-code; ANY matches every protocol. Fixed at creation - it cannot be changed by a later update (required)
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
      --source-ip ip           Match only traffic originating from this IP address or CIDR block (must match --ip-version). Leave unset to allow any source IP
      --source-mac string      Match only traffic originating from this MAC address. Format: aa:bb:cc:dd:ee:ff. Leave unset to allow any source MAC
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Allow inbound SSH (TCP port 22) from any source
ionosctl compute firewallrule create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --protocol TCP --direction INGRESS --port-range-start 22 --port-range-end 22 --name "Allow SSH"

# Allow inbound ICMP echo-request (ping, type 8) from a single source IP only
ionosctl compute firewallrule create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --protocol ICMP --direction INGRESS --icmp-type 8 --source-ip 192.0.2.10 --name "Allow ping from admin host"
```

