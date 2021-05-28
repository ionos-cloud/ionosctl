---
description: Create/Reserve an IpBlock
---

# IpblockCreate

## Usage

```text
ionosctl ipblock create [flags]
```

## Aliases

For `ipblock` command:
```text
[ipb]
```

For `create` command:
```text
[c]
```

## Description

Use this command to create/reserve an IpBlock in a specified location that can be used by resources within any Virtual Data Centers provisioned in that same location. An IpBlock consists of one or more static IP addresses. The name, size of the IpBlock can be set.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Location

## Options

```text
  -u, --api-url string     Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [IpBlockId Name Location Size Ips State] (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               help for create
  -l, --location string    Location of the IpBlock (required)
  -n, --name string        Name of the IpBlock
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
      --size int           Size of the IpBlock (default 2)
  -t, --timeout int        Timeout option for Request for IpBlock creation [seconds] (default 60)
  -w, --wait-for-request   Wait for the Request for IpBlock creation to be executed
```

## Examples

```text
ionosctl ipblock create --ipblock-name NAME --ipblock-location LOCATION_ID --ipblock-size IPBLOCK_SIZE
```

