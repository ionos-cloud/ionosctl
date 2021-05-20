---
description: List Kubernetes Versions
---

# K8sVersionList

## Usage

```text
ionosctl k8s version list [flags]
```

## Description

Use this command to retrieve all available Kubernetes versions.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl k8s version list 
[1.18.16 1.18.15 1.18.12 1.18.5 1.18.9 1.19.8]
```

