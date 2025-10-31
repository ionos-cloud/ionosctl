---
description: "List Kubernetes NodePools"
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
* filter by property: [name datacenterId nodeCount cpuFamily serverType coresCount ramSize availabilityZone storageType storageSize k8sVersion maintenanceWindow autoScaling labels annotations publicIps availableUpgradeVersions]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* K8s Cluster Id

## Options

```text
  -a, --all                 List all resources without the need of specifying parent ID name.
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string   The unique K8s Cluster Id (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [NodePoolId Name K8sVersion DatacenterId NodeCount CpuFamily ServerType StorageType State LanIds CoresCount RamSize AvailabilityZone StorageSize MaintenanceWindow AutoScaling PublicIps AvailableUpgradeVersions Annotations Labels ClusterId] (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           pagination limit: Maximum number of items to return per request (default 50)
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
      --offset int          pagination offset: Number of items to skip before starting to collect the results
      --order-by string     Limits results to those containing a matching value for a specific property
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl k8s nodepool list --cluster-id CLUSTER_ID
```

