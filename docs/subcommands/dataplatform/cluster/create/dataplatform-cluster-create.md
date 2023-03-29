---
description: Create Dataplatform Cluster
---

# DataplatformClusterCreate

## Usage

```text
ionosctl dataplatform cluster create [flags]
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

For `create` command:

```text
[c]
```

## Description

The cluster will be provisioned in the datacenter matching the provided datacenterID. Therefore the datacenter must be created upfront and must be accessible by the user issuing the request

## Options

```text
  -i, --datacenter-id string      The ID of the connected datacenter
      --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur (required)
      --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59 (required)
  -n, --name string               The name of your cluster
  -t, --timeout int               Timeout option for Request [seconds] (default 60)
      --version string            The version of your cluster (default "22.11")
  -w, --wait-for-request          Wait for the Request to be executed
```

