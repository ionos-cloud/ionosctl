---
description: "Create a Kubernetes Cluster"
---

# K8sClusterCreate

## Usage

```text
ionosctl compute k8s cluster create [flags]
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

Use this command to create a new Managed Kubernetes Cluster.
Regarding the name for the Kubernetes Cluster, the limit is 63 characters following the rule to begin and end with an alphanumeric character with dashes, underscores, dots and alphanumerics between.
Regarding the Kubernetes Version for the Cluster, if not set via flag, it will be used the default one: `ionosctl compute k8s version get`.
You can wait for the Cluster to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

## Options

```text
      --api-subnets strings     Access to the K8s API server is restricted to these CIDRs. Cluster-internal traffic is not affected by this restriction. If no allowlist is specified, access is not restricted. If an IP without subnet mask is provided, the default value will be used: 32 for IPv4 and 128 for IPv6
  -u, --api-url string          Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [ClusterId Name K8sVersion State MaintenanceWindow Public Location NatGatewayIp NodeSubnet AvailableUpgradeVersions ViableNodePoolVersions S3Bucket ApiSubnetAllowList] (default [ClusterId,Name,K8sVersion,State,MaintenanceWindow,Public,Location])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int               Level of detail for response objects (default 1)
  -F, --filters strings         Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                   Force command to execute without user input
  -h, --help                    Print usage
      --k8s-version string      The K8s version for the Cluster. If not set, the default one will be used
      --limit int               Maximum number of items to return per request (default 50)
      --location string         This attribute is mandatory if the cluster is private. The location must be enabled for your contract, or you must have a data center at that location. This property is not adjustable. Location de/fra/2 is currently unavailable. (default "us/las")
  -n, --name string             The name for the K8s Cluster (default "UnnamedCluster")
      --nat-gateway-ip string   A reserved IP in the given location if using a private cluster. This is the nat gateway IP of the cluster if the cluster is private. This property is immutable. Must be a reserved IP in the same location as the cluster's location. This attribute is mandatory if the cluster is private
      --no-headers              Don't print table headers when table output is used
      --node-subnet string      The node subnet of the cluster, if the cluster is private. This property is optional and immutable. Must be a valid CIDR notation for an IPv4 network prefix of 16 bits length
      --offset int              Number of items to skip before starting to collect the results
      --order-by string         Property to order the results by
  -o, --output string           Desired output format [text|json|api-json] (default "text")
      --public                  The indicator whether the cluster is public or private (default true)
      --query string            JMESPath query string to filter the output
  -q, --quiet                   Quiet output
      --s3bucket string         S3 Bucket name configured for K8s usage
  -t, --timeout int             Timeout option for waiting for Cluster/Request [seconds] (default 600)
  -v, --verbose count           Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request        Wait for the Request for Cluster creation to be executed
  -W, --wait-for-state          Wait for the new Cluster to be in ACTIVE state
```

## Examples

```text
ionosctl compute k8s cluster create --name NAME
```

