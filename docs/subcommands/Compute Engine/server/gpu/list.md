---
description: "List Gpus from a Server"
---

# ServerGpuList

## Usage

```text
ionosctl compute server gpu list [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `gpu` command:

```text
[gpus]
```

For `list` command:

```text
[l ls]
```

## Description

List Gpus from a Server

Use this command to retrieve a list of Gpus attached to a Server.

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
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl compute server gpu list --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

