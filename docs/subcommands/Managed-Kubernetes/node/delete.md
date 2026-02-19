---
description: "Delete a Kubernetes Node"
---

# K8sNodeDelete

## Usage

```text
ionosctl compute k8s node delete [flags]
```

## Aliases

For `node` command:

```text
[n]
```

For `delete` command:

```text
[d]
```

## Description

This command deletes a Kubernetes Node within an existing Kubernetes NodePool in a Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id

## Options

```text
  -a, --all                  Delete all the Kubernetes Nodes within an existing Kubernetes NodePool in a Cluster.
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [NodeId Name K8sVersion PublicIP PrivateIP State] (default [NodeId,Name,K8sVersion,PublicIP,PrivateIP,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --node-id string       The unique K8s Node Id (required)
      --nodepool-id string   The unique K8s Node Pool Id (required)
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl compute k8s node delete --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID
```

