---
description: Kubernetes Kubeconfig Operations
---

# K8sKubeconfig

## Usage

```text
ionosctl k8s-kubeconfig [command]
```

## Description

The sub-command of `ionosctl k8s-kubeconfig` allows you to get the configuration file of a Kubernetes Cluster.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force            Force command to execute without user input
  -h, --help             help for k8s-kubeconfig
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl k8s-kubeconfig get](get.md) | Get the kubeconfig file for a Kubernetes Cluster |

