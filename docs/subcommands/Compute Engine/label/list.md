---
description: "List Labels from Resources"
---

# LabelList

## Usage

```text
ionosctl compute label list [flags]
```

## Aliases

For `list` command:

```text
[l ls]
```

## Description

Use this command to list Labels.

Without --resource-type it lists every label across ALL resources under your account (each row shows key, value, resource type, resource id and URN). To scope to one resource, pass --resource-type plus that resource's id flag(s), e.g. --resource-type server needs --datacenter-id and --server-id.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [key value resourceId resourceType resourceHref]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [URN Key Value ResourceType ResourceId]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The Data Center Id. Required with --resource-type datacenter, server and volume
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --image-id string        The Image Id (private images only). Used with --resource-type image
      --ipblock-id string      The IpBlock Id. Used with --resource-type ipblock
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --resource-type string   Scope the listing to one resource kind (also pass that kind's id flag(s)). If not given, labels across all resources are listed. Can be one of: datacenter, volume, server, snapshot, ipblock, image (required)
      --server-id string       The Server Id (also needs --datacenter-id). Used with --resource-type server
      --snapshot-id string     The Snapshot Id. Used with --resource-type snapshot
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
      --volume-id string       The Volume Id (also needs --datacenter-id). Used with --resource-type volume
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# All labels on the account
ionosctl compute label list

# Labels on one datacenter
ionosctl compute label list --resource-type datacenter --datacenter-id DATACENTER_ID
```

