---
description: "Get a Data Center"
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
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [DatacenterId Name Location State Description Version Features CpuFamily SecAuthProtection IPv6CidrBlock] (default [DatacenterId,Name,Location,CpuFamily,IPv6CidrBlock,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -i, --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl datacenter get --datacenter-id DATACENTER_ID
```

