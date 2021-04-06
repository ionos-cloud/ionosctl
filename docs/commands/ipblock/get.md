---
description: Get an IpBlock
---

# Get

## Usage

```text
ionosctl ipblock get [flags]
```

## Description

Use this command to get information about a specified IpBlock.

Required values to run command:

* IpBlock Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings        Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                help for get
      --ignore-stdin        Force command to execute without user input
      --ipblock-id string   The unique IpBlock Id [Required flag]
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
```

## Examples

```text
ionosctl ipblock get --ipblock-id 3bb77993-dd2a-4845-8115-5001ae87d5e4 
IpBlockId                              Name   Location   Size   Ips                                State
3bb77993-dd2a-4845-8115-5001ae87d5e4   test   us/las     2      [x.x.x.x x.x.x.x]   AVAILABLE
```

