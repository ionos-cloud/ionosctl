---
description: Kubernetes Cluster Operations
---

# K8sCluster

## Usage

```text
ionosctl k8s-cluster [command]
```

## Description

The sub-commands of `ionosctl k8s-cluster` allow you to list, get, create, update, delete Kubernetes Clusters.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [ClusterId,Name,K8sVersion,State,MaintenanceWindow])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force            Force command to execute without user input
  -h, --help             help for k8s-cluster
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl k8s-cluster create](create.md) | Create a Kubernetes Cluster |
| [ionosctl k8s-cluster delete](delete.md) | Delete a Kubernetes Cluster |
| [ionosctl k8s-cluster get](get.md) | Get a Kubernetes Cluster |
| [ionosctl k8s-cluster list](list.md) | List Kubernetes Clusters |
| [ionosctl k8s-cluster update](update.md) | Update a Kubernetes Cluster |

