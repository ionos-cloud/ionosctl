---
description: Get information about the Contract Resources on your account
---

# ContractGet

## Usage

```text
ionosctl contract get [flags]
```

## Aliases

For `contract` command:

```text
[c]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about the Contract Resources on your account. Use `--resource-limits` flag to see specific Contract Resources Limits.

## Options

```text
  -u, --api-url string           Override default host url (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [ContractNumber Owner Status RegistrationDomain CoresPerServer CoresPerContract CoresProvisioned RamPerServer RamPerContract RamProvisioned HddLimitPerVolume HddLimitPerContract HddVolumeProvisioned SsdLimitPerVolume SsdLimitPerContract SsdVolumeProvisioned DasVolumeProvisioned ReservableIps ReservedIpsOnContract ReservedIpsInUse K8sClusterLimitTotal K8sClustersProvisioned NlbLimitTotal NlbProvisioned NatGatewayLimitTotal NatGatewayProvisioned] (default [ContractNumber,Owner,Status,RegistrationDomain])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                    Force command to execute without user input
  -h, --help                     help for get
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --resource-limits string   Specify Resource Limits to see details about it
  -v, --verbose                  see step by step process when running a command
```

## Examples

```text
ionosctl contract get --resource-limits [ CORES|RAM|HDD|SSD|IPS|K8S ]
```

