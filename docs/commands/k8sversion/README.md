---
description: K8s Version Operations
---

# K8sVersion

## Usage

```text
ionosctl k8s-version [command]
```

## Description

The sub-commands of `ionosctl k8s-version` allows you to get the Kubeconfig file of a Cluster.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for k8s-version
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl k8s-version get](get.md) | Get K8s Default Version |
| [ionosctl k8s-version list](list.md) | List K8s Versions |

