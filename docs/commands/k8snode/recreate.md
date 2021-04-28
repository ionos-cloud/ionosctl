---
description: Recreate a K8s Node
---

# Recreate

## Usage

```text
ionosctl k8s-node recreate [flags]
```

## Description

You can recreate a single Kubernetes Node.

Managed Kubernetes starts a process which based on the NodePool's template creates & configures a new Node, waits for status "ACTIVE", and migrates all the Pods from the faulty Node, deleting it once empty. While this operation occurs, the NodePool will have an extra billable "ACTIVE" Node.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string    The unique K8s Cluster Id [Required flag]
      --cols strings         Columns to be printed in the standard output (default [NodeId,Name,K8sVersion,PublicIP,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                 help for recreate
      --ignore-stdin         Force command to execute without user input
      --node-id string       The unique K8s Node Id [Required flag]
      --nodepool-id string   The unique K8s Node Pool Id [Required flag]
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
```

## Examples

```text
ionosctl k8s-node recreate --cluster-id ba5e2960-4068-4aee-b972-092c254769a8 --nodepool-id a274bc0e-efa5-41c0-828d-39e38f4ad361 --node-id 60ef2bd6-0f63-4006-b448-e8e060edba7d 
Warning: Are you sure you want to recreate k8s node (y/N) ? 
y
Status: Command node recreate has been successfully executed
```

