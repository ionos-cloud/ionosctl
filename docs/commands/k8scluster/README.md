---
description: K8s Cluster Operations
---

# K8sCluster

## Usage

```text
ionosctl k8s-cluster [command]
```

## Description

The sub-commands of `ionosctl k8s-cluster` allow you to list, get, create, update, delete K8s Clusters.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [ClusterId,Name,K8sVersion,State,MaintenanceWindow])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for k8s-cluster
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl k8s-cluster create](create.md) | Create a K8s Cluster |
| [ionosctl k8s-cluster delete](delete.md) | Delete a K8s Cluster |
| [ionosctl k8s-cluster get](get.md) | Get a K8s Cluster |
| [ionosctl k8s-cluster list](list.md) | List K8s Clusters |
| [ionosctl k8s-cluster update](update.md) | Update a K8s Cluster |

