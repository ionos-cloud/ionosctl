---
description: K8s Node Operations
---

# K8sNode

## Usage

```text
ionosctl k8s-node [command]
```

## Description

The sub-commands of `ionosctl k8s-node` allow you to list, get, recreate, delete K8s Nodes.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [NodeId,Name,K8sVersion,PublicIP,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for k8s-node
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl k8s-node delete](delete.md) | Delete a K8s Node |
| [ionosctl k8s-node get](get.md) | Get a K8s Node |
| [ionosctl k8s-node list](list.md) | List K8s Nodes |
| [ionosctl k8s-node recreate](recreate.md) | Recreate a K8s Node |

