---
description: "Get a Kubernetes Cluster"
---

# K8sClusterGet

## Usage

```text
ionosctl k8s cluster get [flags]
```

## Aliases

For `cluster` command:

```text
[c]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific Kubernetes Cluster.You can wait for the Cluster to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique K8s Cluster Id (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name K8sVersion State MaintenanceWindow Public Location NatGatewayIp NodeSubnet AvailableUpgradeVersions ViableNodePoolVersions S3Bucket ApiSubnetAllowList] (default [ClusterId,Name,K8sVersion,State,MaintenanceWindow,Public,Location])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout option for waiting for Cluster to be in ACTIVE state [seconds] (default 600)
  -v, --verbose count       Print step-by-step process when running command
  -W, --wait-for-state      Wait for specified Cluster to be in ACTIVE state
```

## Examples

```text
ionosctl k8s cluster get --cluster-id CLUSTER_ID
```

