---
description: "Remove a Kubernetes NodePool LAN"
---

# K8sNodepoolLanRemove

## Usage

```text
ionosctl compute k8s nodepool lan remove [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `remove` command:

```text
[r]
```

## Description

Detach a LAN from a node pool's worker Nodes. The Nodes lose their interface on
that LAN and any routes defined for it. The LAN itself is not deleted. Pass
--all to detach every LAN from the pool.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id

## Options

```text
  -a, --all                  Detach every LAN currently attached to the node pool
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [NodePoolId Name K8sVersion NodeCount DatacenterId State CpuFamily ServerType StorageType LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions Annotations Labels ClusterId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -i, --lan-id int           ID of the LAN to detach from the pool's worker Nodes (required)
      --limit int            Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --nodepool-id string   The unique K8s Node Pool Id (required)
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute k8s nodepool lan remove --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID
```

