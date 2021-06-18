---
description: Get a Resource Share from a Group
---

# ShareGet

## Usage

```text
ionosctl share get [flags]
```

## Aliases

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve the details of a specific Shared Resource available to a specified Group.

Required values to run command:

* Group Id
* Resource Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ShareId EditPrivilege SharePrivilege Type] (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
      --group-id string      The unique Group Id (required)
  -h, --help                 help for get
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -i, --resource-id string   The unique Resource Id (required)
```

## Examples

```text
ionosctl share get --group-id GROUP_ID --resource-id RESOURCE_ID
```

