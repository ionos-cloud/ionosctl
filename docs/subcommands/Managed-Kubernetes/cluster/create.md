---
description: "Create a Kubernetes Cluster (control plane)"
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

Create a new Managed Kubernetes cluster - the control plane only. The cluster
starts with no worker capacity; add worker Nodes afterwards with
`ionosctl compute k8s nodepool create` once the cluster reaches ACTIVE.

Name: up to 63 characters, must begin and end with an alphanumeric character,
with dashes, underscores, dots and alphanumerics in between.

Kubernetes version: if --k8s-version is not set, the current default is used
(see `ionosctl compute k8s version get`). The cluster version is the upper
bound for every node pool's version, so pick it deliberately.

Public vs private:
  * Public (default): the Kubernetes API server is reachable on a public
    endpoint. You may still restrict access with --api-subnets.
  * Private (--public=false): the control plane stays on your network. This
    requires --location and --nat-gateway-ip (a reserved IP in that location),
    and optionally --node-subnet. These private-only properties are immutable.

Use `--wait` (`-w`) to block until the cluster reaches the AVAILABLE state.

## Options

```text
      --api-subnets strings                            Restrict access to the Kubernetes API server to these CIDRs (allow list). Cluster-internal traffic is not affected. If left empty, access is not restricted. An IP given without a subnet mask defaults to /32 for IPv4 and /128 for IPv6. e.g. --api-subnets 203.0.113.0/24,198.51.100.7
  -u, --api-url string                                 Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings                                   Set of columns to be printed on output 
                                                       Available columns: [ClusterId Name K8sVersion State MaintenanceWindow Public Location NatGatewayIp NodeSubnet AvailableUpgradeVersions ViableNodePoolVersions S3Bucket ApiSubnetAllowList]
  -c, --config string                                  Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                                      Level of detail for response objects (default 1)
  -F, --filters strings                                Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                                          Force command to execute without user input
  -h, --help                                           Print usage
      --k8s-version ionosctl compute k8s version get   Kubernetes control-plane version, e.g. 1.29.5. This is the upper bound for every node pool's version. If not set, the current default is used (see ionosctl compute k8s version get)
      --limit int                                      Maximum number of items to return per request (default 50)
      --location string                                Location of a private cluster (mandatory when --public=false, ignored for public clusters). Must be enabled for your contract or host a Data Center you own. Immutable. Location de/fra/2 is currently unavailable (default "us/las")
  -n, --name string                                    Name of the cluster. Up to 63 characters; must start and end with an alphanumeric character, dashes, underscores, dots and alphanumerics in between (default "UnnamedCluster")
      --nat-gateway-ip string                          NAT gateway IP for a private cluster (mandatory when --public=false, ignored for public clusters). Must be a reserved IP in the same location as --location. Immutable
      --no-headers                                     Don't print table headers when table output is used
      --node-subnet string                             Node subnet for a private cluster (optional, ignored for public clusters). Must be a valid IPv4 CIDR with a /16 prefix, e.g. 10.0.0.0/16. Immutable
      --offset int                                     Number of items to skip before starting to collect the results
      --order-by string                                Property to order the results by
  -o, --output string                                  Desired output format [text|json|api-json] (default "text")
      --public                                         Whether the cluster is public (API server on a public endpoint, default) or private (--public=false, control plane on your network; requires --location and --nat-gateway-ip) (default true)
      --query string                                   JMESPath query string to filter the output
  -q, --quiet                                          Quiet output
      --s3bucket string                                Name of an existing IONOS Object Storage (S3) bucket that will receive Kubernetes API-server audit logs for this cluster
  -t, --timeout int                                    Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                                  Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                                           Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a public cluster with the default Kubernetes version, then wait for it to become AVAILABLE
ionosctl compute k8s cluster create --name my-cluster --wait

# Create a public cluster on a specific version, restrict API-server access to two CIDRs, and ship audit logs to an S3 bucket
ionosctl compute k8s cluster create --name prod --k8s-version 1.29.5 --api-subnets 203.0.113.0/24,198.51.100.7/32 --s3bucket my-k8s-audit-logs

# Create a private cluster (control plane on your network) with a reserved NAT gateway IP
ionosctl compute k8s cluster create --name private --public=false --location de/txl --nat-gateway-ip 203.0.113.10 --node-subnet 10.0.0.0/16
```

