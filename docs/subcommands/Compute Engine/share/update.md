---
description: "Change a Group's permissions on a shared resource"
---

# ShareUpdate

## Usage

```text
ionosctl compute share update [flags]
```

## Aliases

For `update` command:

```text
[u up]
```

## Description

Change the permission bits (--edit-privilege, --share-privilege) of an existing Share for a (Group, Resource) pair. Use this to promote a read-only share to editable/re-shareable, or to walk those permissions back.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Group Id
* Resource Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ShareId EditPrivilege SharePrivilege Type GroupId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --edit-privilege       Set whether the Group's members may edit (modify) the shared resource. E.g.: --edit-privilege=true, --edit-privilege=false
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
      --group-id string      The unique Group Id (required)
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -i, --resource-id string   The unique Resource Id (required)
      --share-privilege      Set whether the Group's members may re-share this resource with other Groups. E.g.: --share-privilege=true, --share-privilege=false
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Allow the group to re-share the resource
ionosctl compute share update --group-id GROUP_ID --resource-id RESOURCE_ID --share-privilege

# Revoke edit rights but keep the share in place
ionosctl compute share update --group-id GROUP_ID --resource-id RESOURCE_ID --edit-privilege=false
```

