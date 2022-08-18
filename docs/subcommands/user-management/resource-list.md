---
description: List Resources
---

# ResourceList

## Usage

```text
ionosctl resource list [flags]
```

## Aliases

For `resource` command:

```text
[res]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a full list of existing Resources. To sort list by Resource Type, use `ionosctl resource get` command.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [ResourceId Name SecAuthProtection Type State] (default [ResourceId,Name,SecAuthProtection,Type,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32      Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       When using text output, don't print headers
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl resource list
```

