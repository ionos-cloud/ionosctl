---
description: Create/Reserve an IpBlock
---

# IpblockCreate

## Usage

```text
ionosctl ipblock create [flags]
```

## Description

Use this command to create/reserve an IpBlock in a specified location that can be used by resources within any Virtual Data Centers provisioned in that same location. An IpBlock consists of one or more static IP addresses. The name, size of the IpBlock can be set.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* IpBlock Location

## Options

```text
  -u, --api-url string     Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings       Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force              Force command to execute without user input
  -h, --help               help for create
      --location string    Location of the IpBlock (required)
      --name string        Name of the IpBlock
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
      --size int           Size of the IpBlock (default 2)
      --timeout int        Timeout option for Request for IpBlock creation [seconds] (default 60)
      --wait-for-request   Wait for the Request for IpBlock creation to be executed
```

## Examples

```text
ionosctl ipblock create --ipblock-name testing --ipblock-location us/las --ipblock-size 1
IpBlockId                              Name      Location   Size   Ips         State
bf932826-d71b-4759-a7d0-0028261c1e8d   testing   us/las     1      [x.x.x.x]   BUSY
RequestId: a99bd16c-bf7b-4966-8a30-437b5182226b
Status: Command ipblock create has been successfully executed
```

