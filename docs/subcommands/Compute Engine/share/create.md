---
description: "Create a Resource Share for a Group"
---

# ShareCreate

## Usage

```text
ionosctl compute share create [flags]
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
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ShareId EditPrivilege SharePrivilege Type GroupId] (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --edit-privilege       Set the group's permission to edit privileges on resource
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
      --share-privilege      Set the group's permission to share resource
  -t, --timeout int          Timeout option for Request for Resource to be shared through a Group [seconds] (default 60)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request     Wait for the Request for Resource share to executed
```

## Examples

```text
ionosctl compute share create --group-id GROUP_ID --resource-id RESOURCE_ID
```

