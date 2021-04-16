---
description: Contract Operations
---

# Contract

## Usage

```text
ionosctl contract [command]
```

## Description

The sub-command of `ionosctl contract` allows you to see information about Contract for your account.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [ContractNumber,Owner,Status,RegistrationDomain])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for contract
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl contract get](get.md) | Get information about the Contract on your account |

