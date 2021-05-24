---
description: Create a Resource Share for a Group
---

# ShareCreate

## Usage

```text
ionosctl share create [flags]
```

## Description

Use this command to create a specific Resource Share to a Group and optionally allow the setting of permissions for that Resource. As an example, you might use this to grant permissions to use an Image or Snapshot to a specific Group.

Required values to run a command:

* Group Id
* Resource Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ShareId EditPrivilege SharePrivilege Type] (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --edit-privilege       Set the group's permission to edit privileges on resource
  -f, --force                Force command to execute without user input
      --group-id string      The unique Group Id (required)
  -h, --help                 help for create
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -i, --resource-id string   The unique Resource Id (required)
      --share-privilege      Set the group's permission to share resource
  -t, --timeout int          Timeout option for Request for Resource to be shared through a Group [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Resource share to executed
```

## Examples

```text
ionosctl share create --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f --resource-id cefc2175-001f-4b94-8693-6263d731fe8e
ShareId                                EditPrivilege   SharePrivilege
cefc2175-001f-4b94-8693-6263d731fe8e   false           false
RequestId: ffb8e7ba-4a49-4ea5-a97e-e3a61e55c277
Status: Command group add-share has been successfully executed
```

