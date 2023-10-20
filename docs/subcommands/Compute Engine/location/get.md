---
description: "Get a Location"
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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [LocationId Name Features ImageAliases CpuFamily] (default [LocationId,Name,CpuFamily])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -i, --location-id string   The unique Location Id (required)
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl location get --location-id LOCATION_ID
```

