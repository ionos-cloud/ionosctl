---
description: List LANs
---

# LanList

## Usage

```text
ionosctl lan list [flags]
```

## Aliases

For `lan` command:

```text
[l]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of LANs provisioned in a specific Virtual Data Center.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LanId Name Public PccId State] (default [LanId,Name,Public,PccId,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for list
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl lan list --datacenter-id DATACENTER_ID
```

