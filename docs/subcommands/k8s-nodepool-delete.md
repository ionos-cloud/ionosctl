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

## Description

This command deletes a Kubernetes Node Pool within an existing Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string    The unique K8s Cluster Id (required)
  -C, --cols strings         Set of columns to be printed on output 
                             Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily StorageType State CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps PublicIps AvailableUpgradeVersions] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 help for delete
      --nodepool-id string   The unique K8s Node Pool Id (required)
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
```

## Examples

```text
ionosctl k8s nodepool delete --cluster-id ba5e2960-4068-4aee-b972-092c254769a8 --nodepool-id 939811fe-cc13-41e2-8a49-87db58c7a812 
Warning: Are you sure you want to delete k8s node pool (y/N) ? 
y
Status: Command node pool delete has been successfully executed
```

