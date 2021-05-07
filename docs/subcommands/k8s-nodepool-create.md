---
description: Create a Kubernetes NodePool
---

# K8sNodepoolCreate

## Usage

```text
ionosctl k8s nodepool create [flags]
```

## Description

Use this command to create a Node Pool into an existing Kubernetes Cluster. The Kubernetes Cluster must be in state "ACTIVE" before creating a Node Pool. The worker Nodes within the Node Pools will be deployed into an existing Data Center. Regarding the name for the Kubernetes NodePool, the limit is 63 characters following the rule to begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.

You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run a command:

* K8s Cluster Id
* Datacenter Id
* K8s NodePool Name

## Options

```text
  -u, --api-url string            Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string         The unique K8s Cluster Id (required)
      --cols strings              Columns to be printed in the standard output (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores-count int           The total number of cores for the Node (default 2)
      --cpu-family string         CPU Type (default "AMD_OPTERON")
      --datacenter-id string      The unique Data Center Id (required)
      --force                     Force command to execute without user input
  -h, --help                      help for create
      --node-count int            The number of worker Nodes that the Node Pool should contain. Min 1, Max: Determined by the resource availability (default 1)
      --node-zone string          The compute Availability Zone in which the Node should exist (default "AUTO")
      --nodepool-name string      The name for the K8s NodePool (required)
      --nodepool-version string   The K8s version for the NodePool
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
      --ram-size int              The amount of memory for the node in MB, e.g. 2048. Size must be specified in multiples of 1024 MB (1 GB) with a minimum of 2048 MB (default 2048)
      --storage-size int          The total allocated storage capacity of a Node (default 10)
      --storage-type string       Storage Type (default "HDD")
      --timeout int               Timeout option for waiting for NodePool/Request [seconds] (default 600)
      --wait-for-state            Wait for the new NodePool to be in ACTIVE state
```

## Examples

```text
ionosctl k8s nodepool create --datacenter-id 3af92af6-c2eb-41e0-b946-6e7ba321abf2 --cluster-id ba5e2960-4068-4aee-b972-092c254769a8 --nodepool-name test1234
NodePoolId                             Name       K8sVersion   NodeCount   DatacenterId                           State
a274bc0e-efa5-41c0-828d-39e38f4ad361   test1234   1.19.8       2           3af92af6-c2eb-41e0-b946-6e7ba321abf2   DEPLOYING
```

