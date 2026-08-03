---
description: "Update a NAT Gateway Rule"
---

# NatgatewayRuleUpdate

## Usage

```text
ionosctl compute natgateway rule update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Use this command to update the matching criteria or masquerade IP of an existing SNAT rule. Only the flags you pass are changed; the rest keep their current values.

The same constraints as on create apply: `--ip` must be a public IP assigned to the parent gateway, and a target port range is only valid for TCP/UDP (not ICMP).

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state before returning.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway Rule Id

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
      --ip ip                  Public IP used to masquerade matched outbound packets. Must be one of the public IPs already assigned to the parent NAT Gateway
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            New human-friendly name for the rule
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --port-range-end int     Last destination port (inclusive) the rule matches. Only applies to TCP/UDP (default 1)
      --port-range-start int   First destination port (inclusive) the rule matches. Only applies to TCP/UDP (default 1)
  -p, --protocol string        Protocol the rule matches: TCP, UDP, ICMP or ALL. A target port range is only valid for TCP/UDP; with ICMP the target port range must not be set
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -i, --rule-id string         The unique Rule Id (required)
      --source-subnet string   Source subnet (CIDR) matched against each packet's source IP, e.g. 10.0.1.0/24
      --target-subnet string   Destination subnet (CIDR) matched against each packet's destination IP. Match any destination if unset
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a rule
ionosctl compute natgateway rule update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID --name renamed-rule

# Point the rule at a different gateway public IP and widen its source subnet
ionosctl compute natgateway rule update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID --ip 203.0.113.11 --source-subnet 10.0.0.0/16
```

