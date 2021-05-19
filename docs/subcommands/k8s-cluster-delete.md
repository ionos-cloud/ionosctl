---
description: Delete a Kubernetes Cluster
---

# K8sClusterDelete

## Usage

```text
ionosctl k8s cluster delete [flags]
```

## Description

This command deletes a Kubernetes cluster. The cluster cannot contain any NodePools when deleting.

You can wait for Request for the Cluster deletion to be executed using `--wait-for-request` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string   The unique K8s Cluster Id (required)
      --cols strings        Columns to be printed in the standard output (default [ClusterId,Name,K8sVersion,Public,State,MaintenanceWindow])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force               Force command to execute without user input
  -h, --help                help for delete
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
      --timeout int         Timeout option for waiting for Request [seconds] (default 600)
      --wait-for-request    Wait for the Request for Cluster deletion to be executed
```

## Examples

```text
ionosctl k8s cluster delete --cluster-id 01d870e6-4118-4396-90bd-917fda3e948d 
Warning: Are you sure you want to delete K8s cluster (y/N) ? 
y
RequestId: ea736d72-9c49-4c1e-88a5-a15c05329f40
Status: Command k8s cluster delete has been successfully executed
```

