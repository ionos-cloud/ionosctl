---
description: Get an IpBlock
---

# IpblockGet

## Usage

```text
ionosctl ipblock get [flags]
```

## Aliases

For `ipblock` command:
```text
[ipb]
```

For `get` command:
```text
[g]
```

## Description

Use this command to retrieve the attributes of a specific IpBlock.

Required values to run command:

* IpBlock Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [IpBlockId Name Location Size Ips State] (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                help for get
  -i, --ipblock-id string   The unique IpBlock Id (required)
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
```

## Examples

```text
ionosctl ipblock get --ipblock-id 3bb77993-dd2a-4845-8115-5001ae87d5e4 
IpBlockId                              Name   Location   Size   Ips                 State
3bb77993-dd2a-4845-8115-5001ae87d5e4   test   us/las     2      [x.x.x.x x.x.x.x]   AVAILABLE
```

