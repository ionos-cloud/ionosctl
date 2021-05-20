---
description: Update a Resource Share from a Group
---

# ShareUpdate

## Usage

```text
ionosctl share update [flags]
```

## Description

Use this command to update the permissions that a Group has for a specific Resource Share.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Group Id
* Resource Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Columns to be printed in the standard output (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --edit-privilege       Update the group's permission to edit privileges on resource
  -f, --force                Force command to execute without user input
      --group-id string      The unique Group Id (required)
  -h, --help                 help for update
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
      --resource-id string   The unique Resource Id (required)
      --share-privilege      Update the group's permission to share resource
  -t, --timeout int          Timeout option for Request for Resource Share update [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Resource Share update to be executed
```

## Examples

```text
ionosctl share update --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f --resource-id cefc2175-001f-4b94-8693-6263d731fe8e --share-privilege 
ShareId                                EditPrivilege   SharePrivilege
cefc2175-001f-4b94-8693-6263d731fe8e   false           true
RequestId: 0dfccab0-c148-40c8-9794-067d23f79f0e
Status: Command group update-share has been successfully executed
```

