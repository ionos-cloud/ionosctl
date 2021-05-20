---
description: List Resources
---

# ResourceList

## Usage

```text
ionosctl resource list [flags]
```

## Description

Use this command to get a full list of existing Resources. To sort list by Resource Type, use `ionosctl resource get` command.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -F, --format strings   Set of fields to be printed on output (default [ResourceId,Name,SecAuthProtection,Type,State])
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl resource list 
ResourceId                             Name                            SecAuthProtection   Type
cefc2175-001f-4b94-8693-6263d731fe8e                                   false               datacenter
d8922413-05f1-48bb-90ed-c2d407e05b1d   IP_BLOCK_2021-04-20T11:02:52Z   false               ipblock
```

