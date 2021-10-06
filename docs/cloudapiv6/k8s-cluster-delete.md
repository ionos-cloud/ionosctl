---
description: Delete a Kubernetes Cluster
---

# K8sClusterDelete

## Usage

```text
ionosctl k8s cluster delete [flags]
```

## Aliases

For `cluster` command:

```text
[c]
```

For `delete` command:

```text
[d]
```

## Description

This command deletes a Kubernetes cluster. The cluster cannot contain any NodePools when deleting.

You can wait for Request for the Cluster deletion to be executed using `--wait-for-request` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique K8s Cluster Id (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name K8sVersion State MaintenanceWindow AvailableUpgradeVersions ViableNodePoolVersions Public GatewayIp S3Bucket ApiSubnetAllowList] (default [ClusterId,Name,K8sVersion,Public,State,MaintenanceWindow])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout option for waiting for Request [seconds] (default 600)
  -v, --verbose             Print step-by-step process when running command
  -w, --wait-for-request    Wait for the Request for Cluster deletion to be executed
```

## Examples

```text
ionosctl k8s cluster delete --cluster-id CLUSTER_ID
```

