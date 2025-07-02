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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily ServerType StorageType State LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions Annotations Labels ClusterId] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -i, --lan-id int           The unique LAN Id of existing LANs to be detached from worker Nodes (required)
      --no-headers           Don't print table headers when table output is used
      --nodepool-id string   The unique K8s Node Pool Id (required)
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl k8s nodepool lan remove --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID
```

