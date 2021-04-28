---
description: Delete a K8s Cluster
---

# Delete

## Usage

```text
ionosctl k8s-cluster delete [flags]
```

## Description

This command deletes a Kubernetes cluster. The cluster cannot contain any node pools when deleting.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string   The unique K8s Cluster Id [Required flag]
      --cols strings        Columns to be printed in the standard output (default [ClusterId,Name,K8sVersion,State,MaintenanceWindow])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                help for delete
      --ignore-stdin        Force command to execute without user input
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
```

## Examples

```text
ionosctl k8s-cluster delete --cluster-id 01d870e6-4118-4396-90bd-917fda3e948d 
Warning: Are you sure you want to delete K8s cluster (y/N) ? 
y
RequestId: ea736d72-9c49-4c1e-88a5-a15c05329f40
Status: Command k8s-cluster delete has been successfully executed
```

