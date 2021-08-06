---
description: List Kubernetes NodePools
---

# K8sNodepoolList

## Usage

```text
ionosctl k8s nodepool list [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of all contained NodePools in a selected Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cluster-id string   The unique K8s Cluster Id (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily StorageType State LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps PublicIps AvailableUpgradeVersions] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                help for list
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
```

## Examples

```text
ionosctl k8s nodepool list --cluster-id CLUSTER_ID
```

