---
description: "Update a Resource Share from a Group"
---

# ShareUpdate

## Usage

```text
ionosctl share update [flags]
```

## Aliases

For `update` command:

```text
[u up]
```

## Description

Use this command to update the permissions that a Group has for a specific Resource Share.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Group Id
* Resource Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud'|'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ShareId EditPrivilege SharePrivilege Type GroupId] (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --edit-privilege       Update the group's permission to edit privileges on resource
  -f, --force                Force command to execute without user input
      --group-id string      The unique Group Id (required)
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -i, --resource-id string   The unique Resource Id (required)
      --share-privilege      Update the group's permission to share resource
  -t, --timeout int          Timeout option for Request for Resource Share update [seconds] (default 60)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait-for-request     Wait for the Request for Resource Share update to be executed
```

## Examples

```text
ionosctl share update --group-id GROUP_ID --resource-id RESOURCE_ID --share-privilege
```

