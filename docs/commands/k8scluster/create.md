---
description: Create a K8s Cluster
---

# Create

## Usage

```text
ionosctl k8s-cluster create [flags]
```

## Description

Use this command to create a new Managed Kubernetes Cluster.

Required values to run a command:

* K8s Cluster Name

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-name string      The name for the K8s Cluster [Required flag]
      --cluster-version string   The K8s version for the Cluster (default "1.19.8")
      --cols strings             Columns to be printed in the standard output (default [ClusterId,Name,K8sVersion,AvailableUpgradeVersions,ViableNodePoolVersions,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                     help for create
      --ignore-stdin             Force command to execute without user input
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
```

## Examples

```text
ionosctl k8s-cluster create --cluster-name demoTest
ClusterId                              Name       K8sVersion   AvailableUpgradeVersions   ViableNodePoolVersions   State
29d9b0c4-351d-4c9e-87e1-201cc0d49afb   demoTest   1.19.8       []                         []                       DEPLOYING
RequestId: 583ba6ae-dd0b-4c68-8fb2-41b3d7bc471b
Status: Command k8s-cluster create has been successfully executed
```

