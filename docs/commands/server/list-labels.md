---
description: List Labels from a Server
---

# ListLabels

## Usage

```text
ionosctl server list-labels [flags]
```

## Description

Use this command to list all Labels from a specified Server.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
  -h, --help                   help for list-labels
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id [Required flag]
```

## Examples

```text
ionosctl server list-labels --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --server-id 27dde318-f0d4-4f97-a04d-9dafe4a89637 
Key    Value
test   test
```

