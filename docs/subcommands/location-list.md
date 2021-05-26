---
description: List Locations
---

# LocationList

## Usage

```text
ionosctl location list [flags]
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

Use this command to get a list of available locations to create objects on.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [LocationId Name Features ImageAliases] (default [LocationId,Name,Features])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl location list
LocationId   Name        Features
de/fra       frankfurt   [SSD]
us/las       lasvegas    [SSD]
us/ewr       newark      [SSD]
de/txl       berlin      [SSD]
gb/lhr       london      [SSD]
```

