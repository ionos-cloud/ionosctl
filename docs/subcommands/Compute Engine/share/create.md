---
description: "Grant a Group access to a specific resource"
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

Grant a Group access to one specific existing resource, creating a Share for the (Group, Resource) pair. Use this to hand a Group a concrete datacenter, image, snapshot, IP block, etc. - separate from the contract-wide privileges you set on the Group itself.

By default the share grants read/use access only. Add --edit-privilege to let members modify the resource, and/or --share-privilege to let them re-share it with other Groups. Find shareable resource IDs with `ionosctl compute resource list`.

Required values to run a command:

* Group Id
* Resource Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ShareId EditPrivilege SharePrivilege Type GroupId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --edit-privilege       Also allow the Group's members to edit (modify) the shared resource, not just view/use it
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
      --share-privilege      Also allow the Group's members to re-share this resource with other Groups
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Give a group read/use access to a datacenter
ionosctl compute share create --group-id GROUP_ID --resource-id DATACENTER_ID

# Give a group full control: edit the resource and re-share it
ionosctl compute share create --group-id GROUP_ID --resource-id RESOURCE_ID --edit-privilege --share-privilege
```

