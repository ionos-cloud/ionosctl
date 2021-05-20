---
description: List Kubernetes Clusters
---

# K8sClusterList

## Usage

```text
ionosctl k8s cluster list [flags]
```

## Description

Use this command to get a list of existing Kubernetes Clusters.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -F, --format strings   Set of fields to be printed on output (default [ClusterId,Name,K8sVersion,Public,State,MaintenanceWindow])
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl k8s cluster list 
ClusterId                              Name    K8sVersion   State
01d870e6-4118-4396-90bd-917fda3e948d   test    1.19.8       ACTIVE
cb47b98f-b8dd-4108-8ac0-b636e36a161d   test3   1.19.8       ACTIVE
```

