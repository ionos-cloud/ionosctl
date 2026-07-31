---
description: "Add a Kubernetes NodePool LAN"
---

# K8sNodepoolLanAdd

## Usage

```text
ionosctl compute k8s nodepool lan add [flags]
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

Attach an existing LAN to a node pool so its worker Nodes gain an interface on
that LAN. The LAN must already exist in the same Data Center as the pool.

Optionally add routes: --network gives destination CIDRs and --gateway-ip the
gateway each is reached through, so Nodes can route to networks behind that
gateway. The two flags are positional - the Nth --network is paired with the Nth
--gateway-ip, so they must have the same number of entries.

Use `--wait` (`-w`) to block until the node pool reaches the AVAILABLE state.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [LanId Dhcp RoutesNetwork RoutesGatewayIp]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --dhcp                 Whether Nodes obtain an IP on this LAN via DHCP. e.g. --dhcp=true, --dhcp=false (default true)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
      --gateway-ip strings   Gateway IPs (IPv4/IPv6) for the corresponding --network routes. Paired positionally with --network, so it must have the same number of entries
  -h, --help                 Print usage
  -i, --lan-id int           ID of an existing LAN (in the pool's Data Center) to attach to the worker Nodes (required)
      --limit int            Maximum number of items to return per request (default 50)
      --network strings      Destination IPv4/IPv6 CIDRs to route via this LAN. Paired positionally with --gateway-ip, so it must have the same number of entries
      --no-headers           Don't print table headers when table output is used
      --nodepool-id string   The unique K8s Node Pool Id (required)
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Attach LAN 2 to a node pool with DHCP
ionosctl compute k8s nodepool lan add --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id 2

# Attach a LAN with a static route (10.0.0.0/24 via 10.1.5.16)
ionosctl compute k8s nodepool lan add --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id 2 --network 10.0.0.0/24 --gateway-ip 10.1.5.16
```

