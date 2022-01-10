---
description: List IpBlocks
---

# IpblockList

## Usage

```text
ionosctl ipblock list [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list IpBlocks.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [ips location size name]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [IpBlockId Name Location Size Ips State] (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -F, --filters strings   Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -M, --max-results int   The maximum number of elements to return
      --order-by string   Limits results to those containing a matching value for a specific property
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl ipblock list
```

