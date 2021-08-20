---
description: Remove a Kubernetes NodePool LAN
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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily StorageType State LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps PublicIps AvailableUpgradeVersions] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 help for remove
  -i, --lan-id int           The unique LAN Id of existing LANs to be detached from worker Nodes (required)
      --nodepool-id string   The unique K8s Node Pool Id (required)
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              see step by step process when running a command
```

## Examples

```text
ionosctl k8s nodepool lan remove --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID
```

