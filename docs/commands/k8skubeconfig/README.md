---
description: K8s Kubeconfig Operations
---

# K8sKubeconfig

## Usage

```text
ionosctl k8s-kubeconfig [command]
```

## Description

The sub-command of `ionosctl k8s-kubeconfig` allows you to get the Kubeconfig file of a Cluster.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for k8s-kubeconfig
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl k8s-kubeconfig get](get.md) | Get a K8s Cluster Kubeconfig |

