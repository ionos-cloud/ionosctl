---
description: "Get a Resource Share from a Group"
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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ShareId EditPrivilege SharePrivilege Type GroupId] (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                Force command to execute without user input
      --group-id string      The unique Group Id (required)
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -i, --resource-id string   The unique Resource Id (required)
  -t, --timeout duration     Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait                 Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl share get --group-id GROUP_ID --resource-id RESOURCE_ID
```

