---
description: List Kubernetes Nodes
---

# K8sNodeList

## Usage

```text
ionosctl k8s node list [flags]
```

## Aliases

For `node` command:

```text
[n]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of existing Kubernetes Nodes.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [NodeId Name K8sVersion PublicIP PrivateIP State] (default [NodeId,Name,K8sVersion,PublicIP,PrivateIP,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 help for list
      --nodepool-id string   The unique K8s Node Pool Id (required)
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
```

## Examples

```text
ionosctl k8s node list --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID
```

