---
description: Kubernetes NodePool Operations
---

# K8sNodepool

## Usage

```text
ionosctl k8s-nodepool [command]
```

## Description

The sub-commands of `ionosctl k8s-nodepool` allow you to list, get, create, update, delete Kubernetes NodePools.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for k8s-nodepool
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl k8s-nodepool create](create.md) | Create a Kubernetes NodePool |
| [ionosctl k8s-nodepool delete](delete.md) | Delete a Kubernetes NodePool |
| [ionosctl k8s-nodepool get](get.md) | Get a Kubernetes NodePool |
| [ionosctl k8s-nodepool list](list.md) | List Kubernetes NodePools |
| [ionosctl k8s-nodepool update](update.md) | Update a Kubernetes NodePool |

