---
description: Update a Kubernetes Cluster
---

# K8sClusterUpdate

## Usage

```text
ionosctl k8s cluster update [flags]
```

## Aliases

For `cluster` command:

```text
[c]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update the name, Kubernetes version, maintenance day and maintenance time of an existing Kubernetes Cluster.

You can wait for the Cluster to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id

## Options

```text
      --api-subnets strings       Access to the K8s API server is restricted to these CIDRs. Cluster-internal traffic is not affected by this restriction. If no allowlist is specified, access is not restricted. If an IP without subnet mask is provided, the default value will be used: 32 for IPv4 and 128 for IPv6. This will overwrite the existing ones
  -i, --cluster-id string         The unique K8s Cluster Id (required)
  -D, --depth int32               Controls the detail depth of the response objects. Max depth is 10.
      --k8s-version string        The K8s version for the Cluster
      --maintenance-day string    The day of the week for Maintenance Window has the English day format as following: Monday or Saturday
      --maintenance-time string   The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00
  -n, --name string               The name for the K8s Cluster
      --s3bucket string           S3 Bucket name configured for K8s usage. It will overwrite the previous value
  -t, --timeout int               Timeout option for waiting for Cluster to be in ACTIVE state after updating [seconds] (default 600)
  -W, --wait-for-state            Wait for specified Cluster to be in ACTIVE state after updating
```

## Examples

```text
ionosctl k8s cluster update --cluster-id CLUSTER_ID --name NAME
```

