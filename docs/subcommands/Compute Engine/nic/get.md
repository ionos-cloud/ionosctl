---
description: "Get a NIC"
---

# NicGet

## Usage

```text
ionosctl nic get [flags]
```

## Aliases

For `nic` command:

```text
[n]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified NIC from specified Data Center and Server.

Required values to run command:

* Data Center Id
* Server Id
* NIC Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Name Dhcp LanId Ips IPv6Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac DHCPv6 IPv6CidrBlock] (default [NicId,Name,Dhcp,LanId,Ips,IPv6Ips,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
  -i, --nic-id string          The unique NIC Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl nic get --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID
```

