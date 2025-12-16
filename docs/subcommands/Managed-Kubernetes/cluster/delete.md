---
description: "Delete a Kubernetes Cluster"
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
  -a, --all                 Delete all the Kubernetes clusters.
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string   The unique K8s Cluster Id (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name K8sVersion State MaintenanceWindow Public Location NatGatewayIp NodeSubnet AvailableUpgradeVersions ViableNodePoolVersions S3Bucket ApiSubnetAllowList] (default [ClusterId,Name,K8sVersion,State,MaintenanceWindow,Public,Location])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout option for waiting for Request [seconds] (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request    Wait for the Request for Cluster deletion to be executed
```

## Examples

```text
ionosctl k8s cluster delete --cluster-id CLUSTER_ID
```

