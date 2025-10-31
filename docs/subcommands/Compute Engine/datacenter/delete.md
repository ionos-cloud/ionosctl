---
description: "Delete a Data Center"
---

# DatacenterDelete

## Usage

```text
ionosctl datacenter delete [flags]
```

## Aliases

For `datacenter` command:

```text
[d dc vdc]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Virtual Data Center from your account. This will remove all objects within the VDC and remove the VDC object itself.

NOTE: This is a highly destructive operation which should be used with extreme caution!

Required values to run command:

* Data Center Id

## Options

```text
  -a, --all                    Delete all the Datacenters.
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [DatacenterId Name Location State Description Version Features CpuFamily SecAuthProtection IPv6CidrBlock] (default [DatacenterId,Name,Location,CpuFamily,IPv6CidrBlock,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -i, --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for Data Center deletion [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request       Wait for the Request for Data Center deletion
```

## Examples

```text
ionosctl datacenter delete --datacenter-id DATACENTER_ID
ionosctl datacenter delete --datacenter-id DATACENTER_ID --force --wait-for-request
```

