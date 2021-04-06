---
description: List IpBlocks
---

# List

## Usage

```text
ionosctl ipblock list [flags]
```

## Description

Use this command to list IpBlocks.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for list
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl ipblock list 
IpBlockId                              Name   Location   Size   Ips                                State
bf932826-d71b-4759-a7d0-0028261c1e8d   demo   us/las     1      [x.x.x.x]                  AVAILABLE
3bb77993-dd2a-4845-8115-5001ae87d5e4   test   us/las     2      [x.x.x.x x.x.x.x]   AVAILABLE
```

