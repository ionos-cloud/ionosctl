---
description: "List Labels from Resources"
---

# LabelList

## Usage

```text
ionosctl label list [flags]
```

## Aliases

For `list` command:

```text
[l ls]
```

## Description

Use this command to list all Labels from all Resources under your account. If you want to list all Labels from a specific Resource, use `--resource-type` option together with the Resource Id: `--datacenter-id`, `--server-id`, `--volume-id`.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [key value resourceId resourceType resourceHref]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [URN Key Value ResourceType ResourceId] (default [URN,Key,Value,ResourceType,ResourceId])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --image-id string        The unique Image Id(note: only private images supported)
      --ipblock-id string      The unique IpBlock Id
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             Don't print table headers when table output is used
      --order-by string        Limits results to those containing a matching value for a specific property
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --resource-type string   Type of resource to list labels from. Can be one of: datacenter, volume, server, snapshot, ipblock, image (required)
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
  -v, --verbose                Print step-by-step process when running command
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label list

ionosctl label list --resource-type datacenter --datacenter-id DATACENTER_ID
```

