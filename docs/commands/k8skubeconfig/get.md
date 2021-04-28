---
description: Get a K8s Cluster Kubeconfig
---

# Get

## Usage

```text
ionosctl k8s-kubeconfig get [flags]
```

## Description

Use this command to retrieve the kubeconfig file for a given Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string   The unique K8s Cluster Id [Required flag]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                help for get
      --ignore-stdin        Force command to execute without user input
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
```

## Examples

```text
ionosctl k8s-kubeconfig get --cluster-id CLUSTER_ID
```

