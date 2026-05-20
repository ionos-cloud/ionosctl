---
description: "Create/Reserve an IpBlock"
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

Use this command to create/reserve an IpBlock in a specified location that can be used by resources within any Virtual Data Centers provisioned in that same location.
An IpBlock consists of one or more static IP addresses. The name, size of the IpBlock can be set.
Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

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
  -l, --location string   Location of the IpBlock. Location de/fra/2 is currently unavailable. (default "de/txl")
  -n, --name string       Name of the IpBlock. If not set, it will automatically be set
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
      --size int          Size of the IpBlock (default 2)
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute ipblock create --name NAME --location LOCATION_ID --size IPBLOCK_SIZE
```

