---
description: "Get a Token from a Server"
---

# ServerTokenGet

## Usage

```text
ionosctl server token get [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `token` command:

```text
[t]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get the Server's jwToken.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ServerId DatacenterId Name AvailabilityZone Cores RAM CpuFamily VmState State TemplateId Type BootCdromId BootVolumeId NicMultiQueue] (default [ServerId,Name,Type,AvailabilityZone,Cores,RAM,CpuFamily,VmState,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -i, --server-id string       The unique Server Id (required)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl server token get --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

