---
description: Get a Location
---

# LocationGet

## Usage

```text
ionosctl location get [flags]
```

## Aliases

For `location` command:
```text
[loc]
```

## Description

Use this command to get information about a specific Location from a Region.

Required values to run command:

* Location Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [LocationId Name Features ImageAliases] (default [LocationId,Name,Features])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 help for get
      --location-id string   The unique Location Id (required)
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
```

## Examples

```text
ionosctl location get --location-id us/las 
LocationId   Name       Features
us/las       lasvegas   [SSD_STORAGE_ZONING SSD]
```

