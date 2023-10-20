---
description: "List Resources from a Group"
---

# GroupResourceList

## Usage

```text
ionosctl group resource list [flags]
```

## Aliases

For `group` command:

```text
[g]
```

For `resource` command:

```text
[res]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Resources assigned to a Group. To see more details about existing Resources, use `ionosctl resource` commands.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ResourceId Name SecAuthProtection Type State] (default [ResourceId,Name,SecAuthProtection,Type,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
      --group-id string     The unique Group Id (required)
  -h, --help                Print usage
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl group resource list --group-id GROUP_ID
```

