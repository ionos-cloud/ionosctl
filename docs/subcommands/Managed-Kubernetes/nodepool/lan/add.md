---
description: "Add a Kubernetes NodePool LAN"
---

# K8sNodepoolLanAdd

## Usage

```text
ionosctl k8s nodepool lan add [flags]
```

## Aliases

For `nodepool` command:

```text
[np]
```

For `add` command:

```text
[a]
```

## Description

Use this command to add a Node Pool LAN into an existing Node Pool.

You can wait for the Node Pool to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run a command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [LanId Dhcp RoutesNetwork RoutesGatewayIp] (default [LanId,Dhcp,RoutesNetwork,RoutesGatewayIp])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --dhcp                 Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false (default true)
  -f, --force                Force command to execute without user input
      --gateway-ip strings   Slice of IPv4 or IPv6 Gateway IPs for the routes. Must contain same number of arguments as --network flag
  -h, --help                 Print usage
  -i, --lan-id int           The unique LAN Id of existing LANs to be attached to worker Nodes (required)
      --limit int            Pagination limit: Maximum number of items to return per request (default 50)
      --network strings      Slice of IPv4 or IPv6 CIDRs to be routed via the interface. Must contain same number of arguments as --gateway-ip flag
      --no-headers           Don't print table headers when table output is used
      --nodepool-id string   The unique K8s Node Pool Id (required)
      --offset int           Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl k8s nodepool lan add --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID
```

