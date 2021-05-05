---
description: List Kubernetes NodePools
---

# K8sNodepoolList

## Usage

```text
ionosctl k8s nodepool list [flags]
```

## Description

Use this command to get a list of all contained NodePools in a selected Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string   The unique K8s Cluster Id (required)
      --cols strings        Columns to be printed in the standard output (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force               Force command to execute without user input
  -h, --help                help for list
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
```

## Examples

```text
ionosctl k8s nodepool list --cluster-id ba5e2960-4068-4aee-b972-092c254769a8 
NodePoolId                             Name        K8sVersion  NodeCount   DatacenterId                           State
939811fe-cc13-41e2-8a49-87db58c7a812   test12345   1.19.8      2           3af92af6-c2eb-41e0-b946-6e7ba321abf2   UPDATING
```

