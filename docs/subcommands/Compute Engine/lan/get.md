---
description: "Get a LAN"
---

# LanGet

## Usage

```text
ionosctl lan get [flags]
```

## Aliases

For `lan` command:

```text
[l]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information of a given LAN.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LanId Name Public PccId IPv6CidrBlock State DatacenterId] (default [LanId,Name,Public,PccId,IPv6CidrBlock,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -i, --lan-id string          The unique LAN Id (required)
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl lan get --datacenter-id DATACENTER_ID --lan-id LAN_ID
```

