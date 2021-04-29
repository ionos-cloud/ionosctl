---
description: Create/Reserve an IpBlock
---

# Create

## Usage

```text
ionosctl ipblock create [flags]
```

## Description

Use this command to create/reserve an IpBlock in a specified location. The name, size options can be provided.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* IpBlock Location

## Options

```text
  -u, --api-url string            Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings              Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                     Force command to execute without user input
  -h, --help                      help for create
      --ipblock-location string   Location of the IpBlock (required)
      --ipblock-name string       Name of the IpBlock
      --ipblock-size int          Size of the IpBlock (default 2)
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
      --timeout int               Timeout option for the IpBlock to be created [seconds] (default 60)
      --wait                      Wait for the IpBlock to be created
```

## Examples

```text
ionosctl ipblock create --ipblock-name testing --ipblock-location us/las --ipblock-size 1
IpBlockId                              Name      Location   Size   Ips         State
bf932826-d71b-4759-a7d0-0028261c1e8d   testing   us/las     1      [x.x.x.x]   BUSY
RequestId: a99bd16c-bf7b-4966-8a30-437b5182226b
Status: Command ipblock create has been successfully executed
```

