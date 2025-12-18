---
description: "Get information about the Contract Resources on your account"
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
  -u, --api-url string           Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [ContractNumber Owner Status RegistrationDomain CoresPerServer CoresPerContract CoresProvisioned RamPerServer RamPerContract RamProvisioned HddLimitPerVolume HddLimitPerContract HddVolumeProvisioned SsdLimitPerVolume SsdLimitPerContract SsdVolumeProvisioned DasVolumeProvisioned ReservableIps ReservedIpsOnContract ReservedIpsInUse K8sClusterLimitTotal K8sClustersProvisioned NlbLimitTotal NlbProvisioned NatGatewayLimitTotal NatGatewayProvisioned] (default [ContractNumber,Owner,Status,RegistrationDomain])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                Level of detail for response objects (default 1)
  -F, --filters strings          Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --limit int                Maximum number of items to return per request (default 50)
      --no-headers               Don't print table headers when table output is used
      --offset int               Number of items to skip before starting to collect the results
      --order-by string          Property to order the results by
  -o, --output string            Desired output format [text|json|api-json] (default "text")
      --query string             JMESPath query string to filter the output
  -q, --quiet                    Quiet output
      --resource-limits string   Specify Resource Limits to see details about it
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl contract get --resource-limits [ CORES|RAM|HDD|SSD|IPS|K8S ]
```

