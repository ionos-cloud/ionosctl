---
description: List Data Centers
---

# DatacenterList

## Usage

```text
ionosctl datacenter list [flags]
```

## Aliases

For `datacenter` command:

```text
[d dc vdc]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a complete list of Virtual Data Centers provisioned under your account. You can setup multiple query parameters.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`
Available Filters:
* filter by property: [description features location name secAuthProtection version]
* filter by metadata: [createdBy createdByUserId createdDate etag lastModifiedBy lastModifiedByUserId lastModifiedDate state]

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [DatacenterId Name Location State Description Version Features CpuFamily SecAuthProtection] (default [DatacenterId,Name,Location,CpuFamily,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          When using text output, don't print headers
      --order-by string     Limits results to those containing a matching value for a specific property
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl datacenter list
ionosctl datacenter list --cols "DatacenterId,Name,Location,Version"
```

