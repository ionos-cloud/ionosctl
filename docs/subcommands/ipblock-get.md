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
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [IpBlockId Name Location Size Ips State] (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                help for get
  -i, --ipblock-id string   The unique IpBlock Id (required)
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             see step by step process when running a command
```

## Examples

```text
ionosctl ipblock get --ipblock-id IPBLOCK_ID
```

