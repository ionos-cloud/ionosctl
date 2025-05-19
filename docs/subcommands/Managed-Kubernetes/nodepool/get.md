---
description: "Get a Kubernetes NodePool"
---

# K8sNodepoolGet

## Usage

```text
ionosctl k8s nodepool get [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific NodePool from an existing Kubernetes Cluster. You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
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
  -t, --timeout int          Timeout option for waiting for NodePool to be in ACTIVE state [seconds] (default 600)
  -v, --verbose              Print step-by-step process when running command
  -W, --wait-for-state       Wait for specified NodePool to be in ACTIVE state
```

## Examples

```text
ionosctl k8s nodepool get --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID
```

