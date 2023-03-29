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

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name publicIP privateIP k8sVersion]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
      --cluster-id string    The unique K8s Cluster Id (required)
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings      Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32    The maximum number of elements to return
      --no-headers           When using text output, don't print headers
      --nodepool-id string   The unique K8s Node Pool Id (required)
      --order-by string      Limits results to those containing a matching value for a specific property
```

## Examples

```text
ionosctl k8s node list --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID
```

