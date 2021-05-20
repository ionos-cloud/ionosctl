---
description: Get Kubernetes Default Version
---

# K8sVersionGet

## Usage

```text
ionosctl k8s version get [flags]
```

## Description

Use this command to retrieve the current default Kubernetes version for Clusters and NodePools.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for get
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl k8s version get 
1.19.8
```

