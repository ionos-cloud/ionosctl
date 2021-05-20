---
description: List Resources Shares through a Group
---

# ShareList

## Usage

```text
ionosctl share list [flags]
```

## Description

Use this command to get a full list of all the Resources that are shared through a specified Group.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ShareId EditPrivilege SharePrivilege Type] (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
      --group-id string   The unique Group Id (required)
  -h, --help              help for list
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
```

## Examples

```text
ionosctl share list --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f 
ShareId                                EditPrivilege   SharePrivilege
cefc2175-001f-4b94-8693-6263d731fe8e   false           false
```

