---
description: "Delete a Kubernetes NodePool"
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
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily ServerType StorageType State LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions Annotations Labels ClusterId] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -i, --nodepool-id string   The unique K8s Node Pool Id (required)
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl k8s nodepool delete --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID
```

