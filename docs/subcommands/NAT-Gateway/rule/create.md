---
description: "Create a NAT Gateway Rule"
---

# NatgatewayRuleCreate

## Usage

```text
ionosctl compute natgateway rule create [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `rule` command:

```text
[r]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a source-NAT (SNAT) rule on a NAT Gateway. The rule masquerades outbound packets whose source matches `--source-subnet` (and, optionally, whose destination matches `--target-subnet` / a target port range) behind the public IP given in `--ip`.

The `--ip` value must be one of the public IPs already assigned to the parent NAT Gateway (`natgateway create/update --ips`); an address not on the gateway is rejected.

Protocol / port constraints: a target port range (`--port-range-start` / `--port-range-end`) is only meaningful for TCP and UDP. If `--protocol` is ICMP the target port range cannot be set. With the default ALL, leave the port range at its default.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state before returning.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Public IP
* Source Subnet

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayRuleId Name Type Protocol SourceSubnet PublicIp TargetSubnet TargetPortRangeStart TargetPortRangeEnd State]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ip ip                  Public IP used to masquerade the source address of matched outbound packets. Must be one of the public IPs already assigned to the parent NAT Gateway (required)
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            Human-friendly name for the rule (default "Unnamed Rule")
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --port-range-end int     Last destination port (inclusive) the rule matches. Only applies to TCP/UDP; ignored for ICMP/ALL (default 1)
      --port-range-start int   First destination port (inclusive) the rule matches. Only applies to TCP/UDP; ignored for ICMP/ALL (default 1)
  -p, --protocol string        Protocol the rule matches: TCP, UDP, ICMP or ALL (default ALL matches every protocol). A target port range is only valid for TCP/UDP; with ICMP the target port range must not be set (required) (default "ALL")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --source-subnet string   Source subnet (CIDR) the rule applies to, matched against each packet's source IP; typically the CIDR of the private LAN whose servers should get outbound access, e.g. 10.0.1.0/24 (required)
      --target-subnet string   Destination subnet (CIDR) the rule applies to, matched against each packet's destination IP. Leave unset to match any destination
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Masquerade all outbound traffic from the 10.0.1.0/24 LAN behind a gateway public IP
ionosctl compute natgateway rule create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name allow-lan --ip 203.0.113.10 --source-subnet 10.0.1.0/24

# TCP-only rule limited to a destination subnet and HTTPS port range
ionosctl compute natgateway rule create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name https-out --ip 203.0.113.10 --source-subnet 10.0.1.0/24 --target-subnet 198.51.100.0/24 --protocol TCP --port-range-start 443 --port-range-end 443
```

