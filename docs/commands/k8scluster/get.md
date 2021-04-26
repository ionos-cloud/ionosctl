---
description: Get a K8s Cluster
---

# Get

## Usage

```text
ionosctl k8s-cluster get [flags]
```

## Description

Use this command to retrieve details about a specific K8s Cluster.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string   The unique K8s Cluster Id [Required flag]
      --cols strings        Columns to be printed in the standard output (default [ClusterId,Name,K8sVersion,AvailableUpgradeVersions,ViableNodePoolVersions,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                help for get
      --ignore-stdin        Force command to execute without user input
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
```

## Examples

```text
ionosctl k8s-cluster get --cluster-id cb47b98f-b8dd-4108-8ac0-b636e36a161d 
ClusterId                              Name    K8sVersion   AvailableUpgradeVersions   ViableNodePoolVersions                           State
cb47b98f-b8dd-4108-8ac0-b636e36a161d   test3   1.19.8       []                         [1.18.16 1.18.15 1.18.12 1.18.5 1.18.9 1.19.8]   ACTIVE
```

