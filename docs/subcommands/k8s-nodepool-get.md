---
description: Get a Kubernetes NodePool
---

# K8sNodepoolGet

## Usage

```text
ionosctl k8s nodepool get [flags]
```

## Description

Use this command to retrieve details about a specific NodePool from an existing Kubernetes Cluster. You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Columns to be printed in the standard output (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 help for get
      --nodepool-id string   The unique K8s Node Pool Id (required)
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout option for waiting for NodePool to be in ACTIVE state [seconds] (default 600)
  -W, --wait-for-state       Wait for specified NodePool to be in ACTIVE state
```

## Examples

```text
ionosctl k8s nodepool get --cluster-id ba5e2960-4068-4aee-b972-092c254769a8 --nodepool-id 939811fe-cc13-41e2-8a49-87db58c7a812 
NodePoolId                             Name        K8sVersion  NodeCount   DatacenterId                           State
939811fe-cc13-41e2-8a49-87db58c7a812   test12345   1.19.8      2           3af92af6-c2eb-41e0-b946-6e7ba321abf2   UPDATING
```

