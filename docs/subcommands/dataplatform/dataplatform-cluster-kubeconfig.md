---
description: Get a Dataplatform Cluster's Kubeconfig by ID
---

# DataplatformClusterKubeconfig

## Usage

```text
dataplatform cluster kubeconfig [flags]
```

## Aliases

For `dataplatform` command:

```text
[mdp dp stackable managed-dataplatform]
```

For `cluster` command:

```text
[c]
```

For `kubeconfig` command:

```text
[k]
```

## Description

Get a Dataplatform Cluster's Kubeconfig by ID

## Options

```text
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Version MaintenanceWindow DatacenterId State]
  -h, --help                help for kubeconfig
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl dataplatform cluster kubeconfig --cluster-id <cluster-id>
```

