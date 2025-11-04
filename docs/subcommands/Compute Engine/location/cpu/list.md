---
description: "List available CPU Architecture from a Location"
---

# LocationCpuList

## Usage

```text
ionosctl location cpu list [flags]
```

## Aliases

For `location` command:

```text
[loc]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get information about available CPU Architectures from a specific Location.

Required values to run command:

* Location Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [CpuFamily MaxCores MaxRam Vendor] (default [CpuFamily,MaxCores,MaxRam,Vendor])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Pagination limit: Maximum number of items to return per request (default 50)
      --location-id string   The unique Location Id (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl location cpu list --location-id LOCATION_ID
```

