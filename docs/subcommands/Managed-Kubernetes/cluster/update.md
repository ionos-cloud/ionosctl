---
description: "Update a Kubernetes Cluster (control plane)"
---

# K8sClusterUpdate

## Usage

```text
ionosctl compute k8s cluster update [flags]
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

Update an existing cluster's control-plane settings: name, Kubernetes version,
maintenance window (day + time), API subnet allow list, and S3 audit log bucket.

Upgrading the Kubernetes version: --k8s-version may only be raised, and only to
a version the API currently offers as an upgrade target for this cluster (tab
completion lists them). Downgrades are not supported. Upgrade the cluster before
upgrading its node pools, since a node pool may never exceed the cluster version.

Maintenance window: --maintenance-day and --maintenance-time together define a
weekly window during which IONOS may apply patches and minor upgrades. Setting a
new window without both parts leaves the missing part unchanged.

Note: --api-subnets and --s3bucket overwrite the previous values rather than
merging with them.

Use `--wait` (`-w`) to block until the cluster returns to the AVAILABLE state.

Required values to run command:

* K8s Cluster Id

## Options

```text
      --api-subnets strings       Restrict access to the Kubernetes API server to these CIDRs (allow list). Cluster-internal traffic is not affected. If empty, access is not restricted. An IP without a subnet mask defaults to /32 for IPv4 and /128 for IPv6. Overwrites the existing allow list
  -u, --api-url string            Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string         The unique K8s Cluster Id (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name K8sVersion State MaintenanceWindow Public Location NatGatewayIp NodeSubnet AvailableUpgradeVersions ViableNodePoolVersions S3Bucket ApiSubnetAllowList]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --k8s-version string        New Kubernetes control-plane version, e.g. 1.29.5. Upgrade-only (no downgrades) and must be an offered upgrade target (see tab completion). Upgrade the cluster before its node pools
      --limit int                 Maximum number of items to return per request (default 50)
      --maintenance-day string    Day of the week for the maintenance window, in English (e.g. Monday, Saturday). Set together with --maintenance-time to define the weekly window
      --maintenance-time string   Start time of the maintenance window in HH:mm:ss (24-hour) format, e.g. 08:00:00. Set together with --maintenance-day
  -n, --name string               New name for the cluster. Up to 63 characters; must start and end with an alphanumeric character
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --public                    Whether the cluster is public or private. Note: the public/private nature of a cluster is fixed at creation and cannot be changed afterwards (default true)
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --s3bucket string           Name of the IONOS Object Storage (S3) bucket that receives API-server audit logs. Overwrites the previous value
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a cluster
ionosctl compute k8s cluster update --cluster-id CLUSTER_ID --name new-name

# Upgrade the control-plane Kubernetes version and set a Sunday-night maintenance window
ionosctl compute k8s cluster update --cluster-id CLUSTER_ID --k8s-version 1.29.5 --maintenance-day Sunday --maintenance-time 02:00:00
```

