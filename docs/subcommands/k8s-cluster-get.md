---
description: Get a Kubernetes Cluster
---

# K8sClusterGet

## Usage

```text
ionosctl k8s cluster get [flags]
```

## Aliases

For `cluster` command:
```text
[c]
```

For `get` command:
```text
[g]
```

## Description

Use this command to retrieve details about a specific Kubernetes Cluster.You can wait for the Cluster to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
  -i, --cluster-id string   The unique K8s Cluster Id (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name K8sVersion State MaintenanceWindow AvailableUpgradeVersions ViableNodePoolVersions Public GatewayIp] (default [ClusterId,Name,K8sVersion,Public,State,MaintenanceWindow])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                help for get
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout option for waiting for Cluster to be in ACTIVE state [seconds] (default 600)
  -W, --wait-for-state      Wait for specified Cluster to be in ACTIVE state
```

## Examples

```text
ionosctl k8s cluster get --cluster-id cb47b98f-b8dd-4108-8ac0-b636e36a161d 
ClusterId                              Name    K8sVersion   State
cb47b98f-b8dd-4108-8ac0-b636e36a161d   test3   1.19.8       ACTIVE
```

