---
description: Delete a Kubernetes Node
---

# Delete

## Usage

```text
ionosctl k8s-node delete [flags]
```

## Description

This command deletes a Kubernetes Node within an existing Kubernetes NodePool in a Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Columns to be printed in the standard output (default [NodeId,Name,K8sVersion,PublicIP,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                Force command to execute without user input
  -h, --help                 help for delete
      --node-id string       The unique K8s Node Id (required)
      --nodepool-id string   The unique K8s Node Pool Id (required)
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
```

## Examples

```text
ionosctl k8s-node delete --cluster-id ba5e2960-4068-4aee-b972-092c254769a8 --nodepool-id a274bc0e-efa5-41c0-828d-39e38f4ad361 --node-id dd520e26-e347-492f-8121-c9dae0495897 
Warning: Are you sure you want to delete k8s node (y/N) ? 
y
Status: Command node delete has been successfully executed
```

