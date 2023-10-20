---
description: "Create a Resource Share for a Group"
---

# ShareCreate

## Usage

```text
ionosctl share create [flags]
```

## Aliases

For `create` command:

```text
[c]
```

## Description

Use this command to create a specific Resource Share to a Group and optionally allow the setting of permissions for that Resource. As an example, you might use this to grant permissions to use an Image or Snapshot to a specific Group.

Required values to run a command:

* Group Id
* Resource Id

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ShareId EditPrivilege SharePrivilege Type GroupId] (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --edit-privilege       Set the group's permission to edit privileges on resource
  -f, --force                Force command to execute without user input
      --group-id string      The unique Group Id (required)
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -i, --resource-id string   The unique Resource Id (required)
      --share-privilege      Set the group's permission to share resource
  -t, --timeout int          Timeout option for Request for Resource to be shared through a Group [seconds] (default 60)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait-for-request     Wait for the Request for Resource share to executed
```

## Examples

```text
ionosctl share create --group-id GROUP_ID --resource-id RESOURCE_ID
```

