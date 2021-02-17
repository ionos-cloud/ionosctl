---
description: List Locations
---

# List

## Usage

```text
ionosctl location list [flags]
```

## Description

Use this command to get a list of available locations to create objects on.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [LocationId,Name,Features])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
  -h, --help             help for list
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Enable verbose output
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

