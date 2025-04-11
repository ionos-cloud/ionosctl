---
description: "Get all Resources of a Type or a specific Resource Type"
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
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -i, --resource-id string   The ID of the specific Resource to retrieve information about
  -t, --timeout duration     Timeout for waiting for resource to reach desired state (default 1m0s)
      --type string          The specific Type of Resources to retrieve information about (required)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait                 Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl resource get --resource-type ipblock
```

