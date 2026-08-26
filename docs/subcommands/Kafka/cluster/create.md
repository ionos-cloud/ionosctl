---
description: "Create a Kafka cluster"
---

# KafkaClusterCreate

## Usage

```text
ionosctl kafka cluster create [flags]
```

## Aliases

For `cluster` command:

```text
[cl]
```

For `create` command:

```text
[c post]
```

## Description

Create a Kafka cluster: 3 broker nodes attached to one of your private LANs.

Sizing (--size). Every size has 3 brokers; sizes differ in the cores/RAM/storage each broker gets:

  XS   1 core   2 GB    195 GB/broker   (585 GB total)   dev / test
  S    2 cores  4 GB    250 GB/broker   (750 GB total)
  M    2 cores  8 GB    400 GB/broker   (1200 GB total)
  L    4 cores  16 GB   800 GB/broker   (2400 GB total)
  XL   8 cores  32 GB   1500 GB/broker  (4500 GB total)  production

Storage is shared by all topics on the cluster; size for (data-rate x retention x replication-factor) across your topics.

Networking. The cluster lives in the LAN --lan-id inside the datacenter --datacenter-id, both in --location. --broker-addresses assigns a private IP to each broker in CIDR notation (e.g. 10.0.0.1/24) taken from that LAN's subnet, so pass exactly 3 addresses. Clients connect over TLS with mutual authentication; port 9093 is appended automatically (e.g. 10.0.0.1:9093).

The cluster is BUSY while it deploys and becomes AVAILABLE when ready; topics and users can only be created once it is AVAILABLE.

Wiki: https://docs.ionos.com/cloud/data-analytics/kafka/api-howtos/create-kafka

## Options

```text
  -u, --api-url string             Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'kafka' and env var 'IONOS_API_URL' (default "https://kafka.%s.ionos.com")
      --broker-addresses strings   Private IP per broker in CIDR notation (e.g. 10.0.0.1/24), one per broker — pass exactly 3, all from the --lan-id subnet. Port 9093 (TLS) is appended for clients (required)
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [Id Name Version Size DatacenterId LanId BrokerAddresses State StateMessage]
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string       ID of the Virtual Data Center holding the LAN the brokers attach to; must be in --location (required)
  -D, --depth int                  Level of detail for response objects (default 1)
  -F, --filters strings            Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
      --lan-id string              ID of the private LAN (inside --datacenter-id) the brokers attach to; clients reach the cluster over this LAN (required)
      --limit int                  Maximum number of items to return per request (default 50)
  -l, --location string            Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, gb/lhr, gb/bhx, us/ewr, us/las, us/mci, fr/par. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string                Human-readable name for the cluster (required)
      --no-headers                 Don't print table headers when table output is used
      --offset int                 Number of items to skip before starting to collect the results
      --order-by string            Property to order the results by
  -o, --output string              Desired output format [text|json|api-json] (default "text")
      --query string               JMESPath query string to filter the output
  -q, --quiet                      Quiet output
      --size string                Cluster size: sets cores/RAM/storage per broker (all sizes have 3 brokers). XS=1c/2GB, S=2c/4GB, M=2c/8GB, L=4c/16GB, XL=8c/32GB. Can be one of: XS, S, M, L, XL (required)
  -t, --timeout int                Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count              Increase verbosity level [-v, -vv, -vvv]
      --version string             Kafka version to deploy, e.g. 3.9.0 (required)
  -w, --wait                       Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl kafka cl create --name my-cluster --version 3.9.0 --size XS --location de/txl --datacenter-id DATACENTER_ID --lan-id LAN_ID --broker-addresses 10.0.0.1/24,10.0.0.2/24,10.0.0.3/24
```

