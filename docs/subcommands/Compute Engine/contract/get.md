---
description: "Get your contract number, owner, status and resource limits"
---

# ContractGet

## Usage

```text
ionosctl compute contract get [flags]
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

Use this command to view your contract details and the resource limits (quotas) enforced on your account.

By default all limits are shown. Use `--resource-limits` to focus the output on a single resource group; each group also shows the amount already provisioned so you can gauge remaining headroom.

## Options

```text
  -u, --api-url string           Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [ContractNumber Owner Status RegistrationDomain CoresPerServer CoresPerContract CoresProvisioned RamPerServer RamPerContract RamProvisioned HddLimitPerVolume HddLimitPerContract HddVolumeProvisioned SsdLimitPerVolume SsdLimitPerContract SsdVolumeProvisioned DasVolumeProvisioned ReservableIps ReservedIpsOnContract ReservedIpsInUse K8sClusterLimitTotal K8sClustersProvisioned NlbLimitTotal NlbProvisioned NatGatewayLimitTotal NatGatewayProvisioned]
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
      --resource-limits string   Restrict the output to one resource-limit group. One of: CORES (vCPUs), RAM, HDD, SSD, DAS (block storage), IPS (reservable IPs), K8S (Kubernetes clusters), NLB (Network Load Balancers), NAT (NAT Gateways)
  -t, --timeout int              Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                     Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Show full contract info and all resource limits
ionosctl compute contract get

# Focus on core (vCPU) limits only
ionosctl compute contract get --resource-limits CORES
```

