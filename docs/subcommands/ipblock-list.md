---
description: List IpBlocks
---

# IpblockList

## Usage

```text
ionosctl ipblock list [flags]
```

## Aliases

For `ipblock` command:
```text
[ipb]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to list IpBlocks.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [IpBlockId Name Location Size Ips State] (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl ipblock list
```

