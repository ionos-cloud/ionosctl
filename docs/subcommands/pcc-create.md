---
description: Create a Private Cross-Connect
---

# PccCreate

## Usage

```text
ionosctl pcc create [flags]
```

## Description

Use this command to create a Private Cross-Connect. You can specify the name and the description for the Private Cross-Connect.

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [PccId,Name,Description])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                    Force command to execute without user input
  -h, --help                     help for create
  -o, --output string            Desired output format [text|json] (default "text")
      --pcc-description string   The description for the Private Cross-Connect
      --pcc-name string          The name for the Private Cross-Connect
  -q, --quiet                    Quiet output
      --timeout int              Timeout option for Request for Private Cross-Connect creation [seconds] (default 60)
      --wait-for-request         Wait for the Request for Private Cross-Connect creation to be executed
```

## Examples

```text
ionosctl pcc create --pcc-name test --pcc-description "test test" --wait-for-request 
PccId                                  Name   Description
e2337b40-52d9-48d2-bcbc-41c5abc29d11   test   test test
RequestId: 64720266-c6e8-4e78-8e31-6754f006dcb1
Status: Command pcc create & wait have been successfully executed
```

