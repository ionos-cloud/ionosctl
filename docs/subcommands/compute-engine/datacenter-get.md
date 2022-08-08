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
[d dc vdc]
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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [DatacenterId Name Location State Description Version Features CpuFamily SecAuthProtection] (default [DatacenterId,Name,Location,CpuFamily,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -i, --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             When using text output, don't print headers
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl datacenter get --datacenter-id DATACENTER_ID
```

