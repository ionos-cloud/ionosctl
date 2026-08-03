---
description: "Update a FirewallRule"
---

# FirewallruleUpdate

## Usage

```text
ionosctl compute firewallrule update [flags]
```

## Aliases

For `firewallrule` command:

```text
[f fr firewall]
```

For `update` command:

```text
[u up]
```

## Description

Update the matching criteria of an existing Firewall Rule on a NIC. Only the flags you pass are changed; the rest keep their current values.

You can retune the match: --direction, --source-mac, --source-ip, --destination-ip, plus --port-range-start/--port-range-end (for a TCP/UDP rule) or --icmp-type/--icmp-code (for an ICMP rule), and --name.

NOTE: the rule's protocol is fixed at creation and CANNOT be changed here (there is no --protocol flag). To change protocol, delete the rule and create a new one. Editing port flags on an ICMP rule (or ICMP flags on a TCP/UDP rule) has no effect.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Firewall Rule Id

## Options

```text
  -u, --api-url string           Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [FirewallRuleId Name Protocol PortRangeStart PortRangeEnd Direction IPVersion State SourceMac SourceIP DestinationIP IcmpCode IcmpType]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int                Level of detail for response objects (default 1)
      --destination-ip -D        When the NIC has multiple IPs, match only traffic directed to this IP address or CIDR block of the NIC (must match --ip-version). Leave unset to allow any target IP. WARNING: the short-hand flag -D is deprecated
  -d, --direction string         Direction of traffic the rule matches: INGRESS (entering the NIC) or EGRESS (leaving the NIC)
  -F, --filters strings          Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -i, --firewallrule-id string   The unique FirewallRule Id (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --icmp-code int            Only for an ICMP rule: match this ICMP code (0-254). Leave unset to allow all codes. Has no effect on a TCP/UDP/ANY rule
      --icmp-type int            Only for an ICMP rule: match this ICMP type (0-254), e.g. 8 = echo request (ping), 0 = echo reply. Leave unset to allow all types. Has no effect on a TCP/UDP/ANY rule
      --ip-version string        The IP version this rule applies to. If --source-ip/--destination-ip are given it must match their version; if omitted it is deduced from those addresses. Can be one of: IPv4, IPv6 (default "IPv4")
      --limit int                Maximum number of items to return per request (default 50)
  -n, --name string              A human-friendly label for the rule. Has no effect on matching; used only to identify the rule in listings
      --nic-id string            The unique NIC Id (required)
      --no-headers               Don't print table headers when table output is used
      --offset int               Number of items to skip before starting to collect the results
      --order-by string          Property to order the results by
  -o, --output string            Desired output format [text|json|api-json] (default "text")
      --port-range-end int       Only for a TCP/UDP rule: last port of the allowed destination-port range (1-65535, inclusive). Set both --port-range-start and --port-range-end; leave both unset to allow all ports. Has no effect on an ICMP/ANY rule (default 1)
      --port-range-start int     Only for a TCP/UDP rule: first port of the allowed destination-port range (1-65535, inclusive). Set both --port-range-start and --port-range-end; leave both unset to allow all ports. Has no effect on an ICMP/ANY rule (default 1)
      --query string             JMESPath query string to filter the output
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id (required)
      --source-ip ip             Match only traffic originating from this IP address or CIDR block (must match --ip-version). Leave unset to allow any source IP
      --source-mac string        Match only traffic originating from this MAC address. Format: aa:bb:cc:dd:ee:ff. Leave unset to allow any source MAC
  -t, --timeout int              Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                     Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a rule
ionosctl compute firewallrule update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID --name "New name" --wait

# Widen an existing TCP rule to the HTTPS port range and restrict it to one source IP
ionosctl compute firewallrule update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID --port-range-start 443 --port-range-end 443 --source-ip 192.0.2.0/24 --wait
```

