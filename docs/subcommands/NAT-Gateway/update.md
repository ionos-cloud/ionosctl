---
description: "Update a NAT Gateway"
---

# NatgatewayUpdate

## Usage

```text
ionosctl compute natgateway update [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update the name or public IPs of a specified NAT Gateway.

Note that `--ips` REPLACES the whole set of public IPs rather than appending to it, so you must pass every IP the gateway should keep. Removing an IP that a SNAT rule still references will break that rule, so update or delete the affected rules first.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state before returning.

Required values to run command:

* Data Center Id
* NAT Gateway Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayId Name PublicIps State DatacenterId]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            New set of reserved public IP addresses (same location as the datacenter). This OVERWRITES the current set, so include every IP the gateway should keep; dropping an IP referenced by a SNAT rule breaks that rule
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            New human-friendly name for the NAT Gateway
  -i, --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a NAT Gateway
ionosctl compute natgateway update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name renamed-gateway

# Replace the gateway's public IPs (this overwrites the existing set)
ionosctl compute natgateway update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --ips 203.0.113.10,203.0.113.12
```

