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
  -a, --all                  Delete all the Resources Share from a specified Group.
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ShareId EditPrivilege SharePrivilege Type GroupId] (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                Force command to execute without user input
      --group-id string      The unique Group Id (required)
  -h, --help                 Print usage
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -i, --resource-id string   The unique Resource Id (required)
  -t, --timeout int          Timeout option for Request for Resource Share deletion [seconds] (default 60)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait-for-request     Wait for the Request for Resource Share deletion to be executed
```

## Examples

```text
ionosctl share delete --group-id GROUP_ID --resource-id RESOURCE_ID --wait-for-request
```

