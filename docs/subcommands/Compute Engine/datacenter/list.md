---
description: "List Data Centers"
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
* filter by property: [name description location version features secAuthProtection ipv6CidrBlock defaultSecurityGroupId]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

## Options

```text
  -u, --api-url string    Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [DatacenterId Name Location State Description Version Features CpuFamily SecAuthProtection IPv6CidrBlock] (default [DatacenterId,Name,Location,CpuFamily,IPv6CidrBlock,State])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32       Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings   Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Pagination offset: Number of items to skip before starting to collect the results
      --order-by string   Limits results to those containing a matching value for a specific property
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl datacenter list
ionosctl datacenter list --cols "DatacenterId,Name,Location,Version"
```

