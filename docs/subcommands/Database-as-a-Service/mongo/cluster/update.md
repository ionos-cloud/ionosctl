---
description: "Update a MongoDB Cluster"
---

# DbaasMongoClusterUpdate

## Usage

```text
ionosctl dbaas mongo cluster update [flags]
```

## Aliases

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `cluster` command:

```text
[c]
```

For `update` command:

```text
[u]
```

## Description


Use this command to update attributes of a MongoDB Cluster. To specify the cluster to update, use the `--cluster-id` flag and the cluster's unique ID you can get from the list command.

Every cluster can update:
* Maintenance window (day and time). To change any of these, you must specify both together (`--maintenance-day` and `--maintenance-time`).
* The display name of the cluster (`--name`).
* The MongoDB major version (`--version`). This can trigger a major upgrade of the cluster, so be sure to check the compatibility of your applications with the new version. Also see the notes in the [API Documentation](https://docs.ionos.com/cloud/databases/mongodb/api-howtos/modify-cluster-attributes/upgrade-the-mongodb-version).
* The backup storage location (`--backup-location`).

Replicaset clusters can update:
* The number of instances in the replicaset (`--instances`).

For enterprise edition clusters, you can also update:
* The memory for each MongoDB host system (`--ram`)
* The CPU Cores for each MongoDB host system (`--cores`)
* Storage size for each MongoDB instance (`--storage-size`)
* Storage type used for the Database (`--storage-type`)
* The number of shards (`--shards`). This is only possible for sharded clusters and requires a sharded_cluster type.
* The MongoDB Connector for Business Intelligence host and port (`--biconnector`) and whether it is enabled (`--biconnector-enabled`).

Business edition clusters currently cannot update their template size (which defines cores, RAM and storage size) this way. This can be done via DCD or API.

Fields which can only be updated under specific conditions:
* Network connection (CIDR, LAN, Datacenter) can only be updated if the amount of shards or instances changes and must be specified together with the new values. LAN and Datacenter must stay the same but need to be specified.
		

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --backup-location string    Object Storage region where snapshots are stored, e.g. de, eu-south-2, eu-central-3. Defaults to the region nearest the cluster
      --biconnector string        New host:port for the MongoDB BI Connector (SQL access for BI tools). Example: r1.m-abcdefgh1234.mongodb.de-fra.ionos.com:27015
      --biconnector-enabled       Enable or disable the BI Connector. If left unset, its state is unchanged. To explicitly disable it, use --biconnector-enabled=false
      --cidr strings              New comma-separated private IPs (with /24 subnet), one per instance. Editable only when the instance/shard count changes, and must be supplied with --datacenter-id and --lan-id. Unavailable range: 10.233.114.0/24
  -i, --cluster-id string         The unique ID of the cluster to update (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name Edition Type URL Instances Shards Health State MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                 New CPU core count PER instance, e.g. 4. Minimum 1. Enterprise only
      --datacenter-id string      Datacenter of the cluster's connection. Only editable while also changing the instance/shard count; the datacenter and LAN must stay the same but must be re-supplied together with --lan-id and --cidr
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           New replicaset size (one primary + n-1 secondaries). Valid values: 1, 3, 5, 7. Only for replicaset clusters (default 1)
      --lan-id string             Numeric LAN ID of the cluster's connection. Must be supplied together with --datacenter-id and --cidr when changing the connection
      --limit int                 Maximum number of items to return per request (default 50)
      --maintenance-day string    New day of the week for the weekly 4-hour maintenance window, e.g. Saturday. Must be changed together with --maintenance-time
      --maintenance-time string   New start time (UTC, HH:MM:SS) of the weekly 4-hour maintenance window, e.g. 16:30:59. Must be changed together with --maintenance-day
  -n, --name string               New human-readable display name for the cluster
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                New RAM PER instance. Multiple of 1024 MB (1 GB), minimum 2 GB. Accepts a unit, e.g. --ram 4GB. Enterprise only
      --shards int32              New shard count for a sharded-cluster (enterprise only). Valid values: 2-32. Scaling shards DOWN is not supported (default 1)
      --storage-size string       New disk size PER instance, e.g. 200GB. Storage can only be grown, not shrunk; better performance above 100 GB. Enterprise only
      --storage-type string       New disk type: HDD, 'SSD Standard' or 'SSD Premium'. Enterprise only (default "\"SSD Standard\"")
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            Upgrade to this MongoDB major version, e.g. 6.0, 7.0. Patch versions are managed automatically. Downgrades are NOT supported; a major upgrade can affect application compatibility
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a cluster and move its maintenance window (day+time must be given together)
ionosctl dbaas mongo cluster update --cluster-id <cluster-id> --name prod-mongo --maintenance-day Sunday --maintenance-time 03:00:00

# Scale an enterprise cluster vertically (per-instance cores/RAM/storage)
ionosctl dbaas mongo cluster update --cluster-id <cluster-id> --cores 4 --ram 8GB --storage-size 200GB
```

