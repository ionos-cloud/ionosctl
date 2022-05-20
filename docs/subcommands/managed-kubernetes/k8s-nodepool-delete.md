---
description: Delete a Kubernetes NodePool
---

# K8sNodepoolDelete

## Usage

```text
ionosctl k8s nodepool delete [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `delete` command:

```text
[d]
```

## Description

This command deletes a Kubernetes Node Pool within an existing Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
  -a, --all                  Delete all the Kubernetes Node Pools within an existing Kubernetes Nodepools.
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily StorageType State LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions Annotations Labels] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -i, --nodepool-id string   The unique K8s Node Pool Id (required)
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl k8s nodepool delete --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID
```

