---
description: List available CPU Architecture from a Location
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
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [CpuFamily MaxCores MaxRam Vendor] (default [CpuFamily,MaxCores,MaxRam,Vendor])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 help for list
  -i, --location-id string   The unique Location Id (required)
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
```

## Examples

```text
ionosctl location cpu list -i LOCATION_ID
```

