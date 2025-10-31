---
description: "Get the kubeconfig file for a Kubernetes Cluster"
---

# K8sKubeconfigGet

## Usage

```text
ionosctl k8s kubeconfig get [flags]
```

## Aliases

For `kubeconfig` command:

```text
[cfg config]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve the kubeconfig file for a given Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string   The unique K8s Cluster Id (required)
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl k8s kubeconfig get --cluster-id CLUSTER_ID
```

