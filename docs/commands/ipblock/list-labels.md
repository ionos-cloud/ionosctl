---
description: List Labels from a IPBlock
---

# ListLabels

## Usage

```text
ionosctl ipblock list-labels [flags]
```

## Description

Use this command to list all Labels from a specified IPBlock.

Required values to run command:

* IPBlock Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings        Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                help for list-labels
      --ignore-stdin        Force command to execute without user input
      --ipblock-id string   The unique IPBlock Id [Required flag]
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
```

## Examples

```text
ionosctl ipblock list-labels --ipblock-id 379a995b-f285-493e-a56a-f32e1cb6dd06 
Key    Value
test   testipblock
```

