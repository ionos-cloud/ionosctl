---
description: Get information about the Contract Resources on your account
---

# ContractGet

## Usage

```text
ionosctl contract get [flags]
```

## Description

Use this command to get information about the Contract Resources on your account. Use `--resource-limits` flag to see specific Contract Resources Limits.

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                    Force command to execute without user input
  -F, --format strings           Set of fields to be printed on output (default [ContractNumber,Owner,Status,RegistrationDomain])
  -h, --help                     help for get
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --resource-limits string   Specify Resource Limits to see details about it
```

## Examples

```text
ionosctl contract get --resource-limits [ CORES|RAM|HDD|SSD|IPS|K8S ]
```

