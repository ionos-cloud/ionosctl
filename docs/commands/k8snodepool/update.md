---
description: Update a Kubernetes NodePool
---

# Update

## Usage

```text
ionosctl k8s-nodepool update [flags]
```

## Description

Use this command to update the number of worker Nodes, the minimum and maximum number of worker Nodes, the add labels, annotations, to update the maintenance day and time, to attach private LANs to a Node Pool within an existing Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
      --annotation-key string     Annotation key. Must be set together with --annotation-value
      --annotation-value string   Annotation value. Must be set together with --annotation-key
  -u, --api-url string            Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-id string         The unique K8s Cluster Id [Required flag]
      --cols strings              Columns to be printed in the standard output (default [NodePoolId,Name,K8sVersion,NodeCount,DatacenterId,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                      help for update
      --ignore-stdin              Force command to execute without user input
      --label-key string          Label key. Must be set together with --label-value
      --label-value string        Label value. Must be set together with --label-key
      --lan-id int                The unique LAN Id of existing LANs to be attached to worker Nodes
      --maintenance-day string    The day of the week for Maintenance Window has the English day format as following: Monday or Saturday
      --maintenance-time string   The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00
      --max-node-count int        The maximum number of worker Nodes that the managed NodePool can scale out. Should be set together with --min-node-count (default 1)
      --min-node-count int        The minimum number of worker Nodes that the managed NodePool can scale in. Should be set together with --max-node-count (default 1)
      --node-count int            The number of worker Nodes that the NodePool should contain (default 1)
      --nodepool-id string        The unique K8s Node Pool Id [Required flag]
      --nodepool-version string   The K8s version for the NodePool. K8s version downgrade is not supported
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
```

## Examples

```text
ionosctl k8s-nodepool update --cluster-id ba5e2960-4068-4aee-b972-092c254769a8 --nodepool-id f01f4d6c-41a9-47c3-a5a5-f3667cc25265 --node-count=1
Status: Command k8s-nodepool update has been successfully executed
```

