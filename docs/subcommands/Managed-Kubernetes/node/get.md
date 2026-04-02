---
description: "Get a Kubernetes Node"
---

# K8sNodeGet

## Usage

```text
ionosctl compute k8s node get [flags]
```

## Aliases

For `node` command:

```text
[n]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific Kubernetes Node.You can wait for the Node to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id

## Options

```text
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
  -i, --node-id string       The unique K8s Node Id (required)
      --nodepool-id string   The unique K8s Node Pool Id (required)
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout option for waiting for Node to be in ACTIVE state [seconds] (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -W, --wait-for-state       Wait for specified Node to be in ACTIVE state
```

## Examples

```text
ionosctl compute k8s node get --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID
```

