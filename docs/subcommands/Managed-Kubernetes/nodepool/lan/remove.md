---
description: "Remove a Kubernetes NodePool LAN"
---

# K8sNodepoolLanRemove

## Usage

```text
ionosctl k8s nodepool lan remove [flags]
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

This command removes a Kubernetes Node Pool LAN from a Node Pool.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id

## Options

```text
  -a, --all                  Remove all FK8s Nodepool Lans.
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily ServerType StorageType State LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions Annotations Labels ClusterId] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -i, --lan-id int           The unique LAN Id of existing LANs to be detached from worker Nodes (required)
      --limit int            Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --nodepool-id string   The unique K8s Node Pool Id (required)
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl k8s nodepool lan remove --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID
```

