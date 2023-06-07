---
description: List Labels from Resources
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
* filter by property: [key resourceHref resourceId resourceType value]
* filter by metadata: [createdBy createdByUserId createdDate etag lastModifiedBy lastModifiedByUserId lastModifiedDate state]

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [Key Value] (default [Key,Value])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ipblock-id string      The unique IpBlock Id
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             When using text output, don't print headers
      --order-by string        Limits results to those containing a matching value for a specific property
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --resource-type string   Type of resource to list labels from. Can be one of: datacenter, volume, server, snapshot, ipblock (required)
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

