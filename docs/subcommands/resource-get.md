---
description: Get all Resources of a Type or a specific Resource Type
---

# ResourceGet

## Usage

```text
ionosctl resource get [flags]
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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ResourceId Name SecAuthProtection Type State] (default [ResourceId,Name,SecAuthProtection,Type,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 help for get
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -i, --resource-id string   The ID of the specific Resource to retrieve information about
      --type string          The specific Type of Resources to retrieve information about
```

## Examples

```text
ionosctl resource get --resource-type ipblock
```

