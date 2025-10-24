---
description: "List Kubernetes NodePool LANs"
---

# K8sNodepoolLanList

## Usage

```text
ionosctl k8s nodepool lan list [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of all contained NodePool LANs in a selected Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [LanId Dhcp RoutesNetwork RoutesGatewayIp] (default [LanId,Dhcp,RoutesNetwork,RoutesGatewayIp])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -M, --max-results int32    The maximum number of elements to return
      --no-headers           Don't print table headers when table output is used
      --nodepool-id string   The unique K8s Node Pool Id (required)
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl k8s nodepool lan list --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID
```

