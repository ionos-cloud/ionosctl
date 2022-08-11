---
description: Get the kubeconfig file for a Data Platform Cluster
---

# DataplatformKubeconfigGet

## Usage

```text
ionosctl dataplatform kubeconfig get [flags]
```

## Aliases

For `dataplatform` command:

```text
[dp]
```

For `kubeconfig` command:

```text
[cfg config]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve the kubeconfig file for a given Data Platform Cluster.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cluster-id string   The unique ID of the Cluster (required)
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          When using text output, don't print headers
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dataplatform kubeconfig get --cluster-id CLUSTER_ID
```

