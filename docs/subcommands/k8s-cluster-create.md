---
description: Create a Kubernetes Cluster
---

# K8sClusterCreate

## Usage

```text
ionosctl k8s cluster create [flags]
```

## Aliases

For `cluster` command:

```text
[c]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a new Managed Kubernetes Cluster. Regarding the name for the Kubernetes Cluster, the limit is 63 characters following the rule to begin and end with an alphanumeric character with dashes, underscores, dots, and alphanumerics between. Regarding the Kubernetes Version for the Cluster, if not set via flag, it will be used the default one: `ionosctl k8s version get`.

You can wait for the Cluster to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run a command:

* Name

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --bucket-names strings   List of S3 bucket configured for K8s usage
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ClusterId Name K8sVersion State MaintenanceWindow AvailableUpgradeVersions ViableNodePoolVersions Public GatewayIp BucketNames ApiSubnetAllowList] (default [ClusterId,Name,K8sVersion,Public,State,MaintenanceWindow])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                  Force command to execute without user input
      --gateway-ip public      The IP address of the gateway used by the Cluster. This is mandatory when public is set to `false` and should not be provided otherwise
  -h, --help                   help for create
      --k8s-version string     The K8s version for the Cluster. If not set, it will be used the default one
  -n, --name string            The name for the K8s Cluster (required)
  -o, --output string          Desired output format [text|json] (default "text")
      --public                 The indicator if the Cluster is public or private (default true)
  -q, --quiet                  Quiet output
      --subnets strings        Access to the K8s API server is restricted to these CIDRs. Cluster-internal traffic is not affected by this restriction. If no allowlist is specified, access is not restricted. If an IP without subnet mask is provided, the default value will be used: 32 for IPv4 and 128 for IPv6
  -t, --timeout int            Timeout option for waiting for Cluster/Request [seconds] (default 600)
  -w, --wait-for-request       Wait for the Request for Cluster creation to be executed
  -W, --wait-for-state         Wait for the new Cluster to be in ACTIVE state
```

## Examples

```text
ionosctl k8s cluster create --name NAME
```

