---
description: List Labels from a Data Center
---

# ListLabels

## Usage

```text
ionosctl datacenter list-labels [flags]
```

## Description

Use this command to list all Labels from a specified Data Center.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
      --force                  Force command to execute without user input
  -h, --help                   help for list-labels
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl datacenter list-labels --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 
Key    Value
test   testdatacenter
```

