---
description: Delete a Resource Share from a Group
---

# ShareDelete

## Usage

```text
ionosctl share delete [flags]
```

## Aliases

For `delete` command:

```text
[d]
```

## Description

This command deletes a Resource Share from a specified Group.

Required values to run command:

* Resource Id
* Group Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ShareId EditPrivilege SharePrivilege Type] (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
      --group-id string      The unique Group Id (required)
  -h, --help                 help for delete
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -i, --resource-id string   The unique Resource Id (required)
  -t, --timeout int          Timeout option for Request for Resource Share deletion [seconds] (default 60)
  -v, --verbose              see step by step process when running a command
  -w, --wait-for-request     Wait for the Request for Resource Share deletion to be executed
```

## Examples

```text
ionosctl share delete --group-id GROUP_ID --resource-id RESOURCE_ID --wait-for-request
```

