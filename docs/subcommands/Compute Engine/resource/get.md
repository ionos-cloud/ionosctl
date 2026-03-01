---
description: "Get all Resources of a Type or a specific Resource Type"
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

Use this command to get all Resources of a Type or a specific Resource Type using its Type and ID.

Required values to run command:

* Type

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ResourceId Name SecAuthProtection Type State] (default [ResourceId,Name,SecAuthProtection,Type,State])
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
  -i, --resource-id string   The ID of the specific Resource to retrieve information about
      --type string          The specific Type of Resources to retrieve information about (required)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl compute resource get --resource-type ipblock
```

