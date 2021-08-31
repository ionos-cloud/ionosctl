---
description: List Resources Shares through a Group
---

# ShareList

## Usage

```text
ionosctl share list [flags]
```

## Aliases

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a full list of all the Resources that are shared through a specified Group.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ShareId EditPrivilege SharePrivilege Type] (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
      --group-id string   The unique Group Id (required)
  -h, --help              Print usage
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl share list --group-id GROUP_ID
```

