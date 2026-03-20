---
description: "Get a Location"
---

# LocationGet

## Usage

```text
ionosctl compute location get [flags]
```

## Aliases

For `location` command:

```text
[loc]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specific Location from a Region.

Required values to run command:

* Location Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [LocationId Name Features ImageAliases CpuFamily] (default [LocationId,Name,CpuFamily])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -i, --location-id string   The unique Location Id (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl compute location get --location-id LOCATION_ID
```

