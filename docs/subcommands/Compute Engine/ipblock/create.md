---
description: "Create/Reserve an IpBlock"
---

# IpblockCreate

## Usage

```text
ionosctl ipblock create [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create/reserve an IpBlock in a specified location that can be used by resources within any Virtual Data Centers provisioned in that same location. An IpBlock consists of one or more static IP addresses. The name, size of the IpBlock can be set.

You can wait for the Request to be executed using `--wait-for-request` option.

## Options

```text
  -u, --api-url string     Override default host URL. If set, this will be preferred over the config file override. If unset, the default will only be used as a fallback (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [IpBlockId Name Location Size Ips State] (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
  -l, --location string    Location of the IpBlock (default "de/txl")
  -n, --name string        Name of the IpBlock. If not set, it will automatically be set
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
      --size int           Size of the IpBlock (default 2)
  -t, --timeout int        Timeout option for Request for IpBlock creation [seconds] (default 60)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait-for-request   Wait for the Request for IpBlock creation to be executed
```

## Examples

```text
ionosctl ipblock create --name NAME --location LOCATION_ID --size IPBLOCK_SIZE
```

