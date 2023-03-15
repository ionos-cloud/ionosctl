---
description: Update Dataplatform Nodepools
---

# DataplatformNodepoolUpdate

## Usage

```text
dataplatform nodepool update [flags]
```

## Aliases

For `dataplatform` command:

```text
[mdp dp stackable managed-dataplatform]
```

For `nodepool` command:

```text
[np]
```

For `update` command:

```text
[u]
```

## Description

Node pools are the resources that powers the DataPlatformCluster.

The following requests allows to alter the existing resources of the cluster

## Options

```text
  -A, --annotations stringToString   Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE (default [])
      --cluster-id string            The UUID of the cluster the nodepool belongs to
      --cols strings                 Set of columns to be printed on output 
                                     Available columns: [Id Name Nodes Cores CpuFamily Ram Storage MaintenanceWindow State AvailabilityZone Labels Annotations]
  -h, --help                         help for update
  -L, --labels stringToString        Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE (default [])
      --maintenance-day string       Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur (required)
      --maintenance-time string      Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59 (required)
      --no-headers                   When using text output, don't print headers
  -n, --node-count int32             The number of nodes that make up the node pool (required)
  -i, --nodepool-id string           The UUID of the cluster the nodepool belongs to
  -t, --timeout int                  Timeout option for Request [seconds] (default 60)
  -w, --wait-for-request             Wait for the Request to be executed
```

