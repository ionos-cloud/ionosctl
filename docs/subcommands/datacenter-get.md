---
description: Get a Data Center
---

# DatacenterGet

## Usage

```text
ionosctl datacenter get [flags]
```

## Aliases

For `datacenter` command:
```text
[d dc]
```

For `get` command:
```text
[g]
```

## Description

Use this command to retrieve details about a Virtual Data Center by using its ID.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [DatacenterId Name Location State Description Version Features SecAuthProtection] (default [DatacenterId,Name,Location,Features,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -i, --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for get
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl datacenter get --datacenter-id DATACENTER_ID
```

