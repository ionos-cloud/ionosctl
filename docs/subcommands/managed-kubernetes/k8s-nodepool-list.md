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

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name datacenterId nodeCount cpuFamily coresCount ramSize availabilityZone storageType storageSize k8sVersion maintenanceWindow autoScaling labels annotations publicIps availableUpgradeVersions gatewayIp]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cluster-id string   The unique K8s Cluster Id (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily StorageType State LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions GatewayIp] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -F, --filters strings     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -M, --max-results int     The maximum number of elements to return
      --no-headers          When using text output, don't print headers
      --order-by string     Limits results to those containing a matching value for a specific property
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl k8s nodepool list --cluster-id CLUSTER_ID
```

