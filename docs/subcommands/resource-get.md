---
description: Get all Resources of a Type or a specific Resource Type
---

# ResourceGet

## Usage

```text
ionosctl resource get [flags]
```

## Description

Use this command to get all Resources of a Type or a specific Resource Type using its Type and ID.

Required values to run command:

* Resource Type

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [ResourceId,Name,SecAuthProtection,Type,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                  Force command to execute without user input
  -h, --help                   help for get
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --resource-id string     The ID of the specific Resource to retrieve information about
      --resource-type string   The specific Type of Resources to retrieve information about
```

## Examples

```text
ionosctl resource get --resource-type ipblock
ResourceId                             Name                            SecAuthProtection   Type
d8922413-05f1-48bb-90ed-c2d407e05b1d   IP_BLOCK_2021-04-20T11:02:52Z   false               ipblock
```

