---
description: "List resources of one type, or get a single resource"
---

# ResourceGet

## Usage

```text
ionosctl compute resource get [flags]
```

## Aliases

For `resource` command:

```text
[res]
```

For `get` command:

```text
[g]
```

## Description

List all resources of a given type, or - if you also pass --resource-id - fetch one specific resource of that type. Valid types are: datacenter, snapshot, image, ipblock, pcc, backupunit, k8s.

Required values to run command:

* Type

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ResourceId Name SecAuthProtection Type State]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -i, --resource-id string   Optional: the ID of a single resource (of the given --resource-type) to fetch. Omit to list all resources of that type
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
      --type string          The type of resources to retrieve. One of: datacenter, snapshot, image, ipblock, pcc, backupunit, k8s (required)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# List every IP block on the contract
ionosctl compute resource get --resource-type ipblock

# Fetch one specific datacenter (to grab its ID/type before sharing it)
ionosctl compute resource get --resource-type datacenter --resource-id DATACENTER_ID
```

