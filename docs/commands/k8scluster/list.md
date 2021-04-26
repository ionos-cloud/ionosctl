---
description: List K8s Clusters
---

# List

## Usage

```text
ionosctl k8s-cluster list [flags]
```

## Description

Use this command to get a list of existing K8s Clusters.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [ClusterId,Name,K8sVersion,AvailableUpgradeVersions,ViableNodePoolVersions,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for list
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl k8s-cluster list 
ClusterId                              Name    K8sVersion   AvailableUpgradeVersions   ViableNodePoolVersions                           State
01d870e6-4118-4396-90bd-917fda3e948d   test    1.19.8       []                         [1.18.16 1.18.15 1.18.12 1.18.5 1.18.9 1.19.8]   ACTIVE
cb47b98f-b8dd-4108-8ac0-b636e36a161d   test3   1.19.8       []                         [1.18.16 1.18.15 1.18.12 1.18.5 1.18.9 1.19.8]   ACTIVE
```

