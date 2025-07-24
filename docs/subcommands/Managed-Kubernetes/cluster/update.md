---
description: "Update a Kubernetes Cluster"
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
  -u, --api-url string            Override default host URL. Preferred over the config file override 'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string         The unique K8s Cluster Id (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name K8sVersion State MaintenanceWindow Public Location NatGatewayIp NodeSubnet AvailableUpgradeVersions ViableNodePoolVersions S3Bucket ApiSubnetAllowList] (default [ClusterId,Name,K8sVersion,State,MaintenanceWindow,Public,Location])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32               Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --k8s-version string        The K8s version for the Cluster
      --maintenance-day string    The day of the week for Maintenance Window has the English day format as following: Monday or Saturday
      --maintenance-time string   The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00
  -n, --name string               The name for the K8s Cluster
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --public                    The indicator whether the cluster is public or private (default true)
  -q, --quiet                     Quiet output
      --s3bucket string           S3 Bucket name configured for K8s usage. It will overwrite the previous value
  -t, --timeout int               Timeout option for waiting for Cluster to be in ACTIVE state after updating [seconds] (default 600)
  -v, --verbose                   Print step-by-step process when running command
  -W, --wait-for-state            Wait for specified Cluster to be in ACTIVE state after updating
```

## Examples

```text
ionosctl k8s cluster update --cluster-id CLUSTER_ID --name NAME
```

