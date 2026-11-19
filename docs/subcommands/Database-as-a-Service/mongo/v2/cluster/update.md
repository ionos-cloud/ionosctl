---
description: "Update a MongoDB Cluster"
---

# DbaasMongoV2ClusterUpdate

## Usage

```text
ionosctl dbaas mongo-v2 cluster update [flags]
```

## Aliases

For `mongo-v2` command:

```text
[mongodb-v2 mg-v2]
```

For `cluster` command:

```text
[c]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update attributes of a MongoDB Cluster. This command uses a combination of GET and PUT to simulate a PATCH operation.

The edition, the instance template of a BUSINESS cluster, the snapshot location, the datacenter and the LAN are fixed at creation. A BUSINESS cluster can change its --instances; an ENTERPRISE cluster its --instances, --cores, --ram and --storage-size (storage can only grow). On an ENTERPRISE cluster created with sharding, --shards adds shards: shards cannot be removed, and sharding cannot be turned on or off after creation. --bi-connector-enabled toggles the BI Connector of an ENTERPRISE cluster.

When the number of instances changes, pass --cidr in the same command so it keeps listing one address per instance (a sharded cluster always has exactly three).

The version can only be upgraded, to one of the CanUpgradeTo versions of `ionosctl dbaas mongo-v2 version list`. Upgrading asks for confirmation: snapshots taken before the upgrade can no longer be restored in place.

The initial user's credentials stay unchanged unless --password is set. The API never returns the password, so changing --user or --database requires --password too.

Required values to run command:

* Cluster Id

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mongodbv2' and env var 'IONOS_API_URL' (default "https://mongodb.%s.ionos.com")
      --bi-connector-enabled      ENTERPRISE only: enable or disable the MongoDB BI Connector
      --cidr strings              The IP addresses with netmask of the cluster instances, replacing the current ones: one per instance, or exactly three with sharding
  -i, --cluster-id string         The unique ID of the Cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name Edition Version Instances State ConnectionUri Cores Ram StorageSize TemplateId ShardingEnabled Shards BiConnectorEnabled Description SnapshotLocation RetentionDays MaintenanceDay MaintenanceTime LogsEnabled MetricsEnabled DatacenterId LanId Cidr Username Database StatusMessage]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int32               ENTERPRISE only: the number of CPU cores per instance, between 1 and 62
      --database string           Database for the initial MongoDB user. Requires --password. Defaults to the cluster's current database
  -D, --depth int                 Level of detail for response objects (default 1)
      --description string        Human-readable description for the cluster
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           The total number of instances in the cluster (one primary and n-1 secondary), between 1 and 9. With sharding, this is the number of replicas per shard
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: gb/bhx, us/ewr. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --logs-enabled              Forward this cluster's logs to the IONOS Logging Service. Has no effect unless that service is activated on your contract, which incurs additional costs
      --maintenance-day string    Day of the week for the MaintenanceWindow. Must be specified together with --maintenance-time. Can be one of: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday
      --maintenance-time string   Time for the MaintenanceWindow. e.g.: 16:30:59. Must be specified together with --maintenance-day
      --metrics-enabled           Forward this cluster's metrics to the IONOS Monitoring Service. Has no effect unless that service is activated on your contract, which incurs additional costs
  -n, --name string               The friendly name of your cluster
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --password string           New password for the initial MongoDB user: 8-256 characters. When omitted, the user's credentials are left unchanged
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                ENTERPRISE only: the amount of memory per instance, between 6GB and 240GB. e.g. --ram 8, --ram 8GB
      --retention-days int32      How many days snapshots are retained. Minimum: 1, Maximum: 365. Reducing it purges the existing snapshots outside the new retention window
      --shards int32              ENTERPRISE only: the new number of shards, up to 5, for a cluster created with sharding. Shards can only be added: they cannot be removed, and sharding cannot be enabled on an existing cluster
      --storage-size string       ENTERPRISE only: the size of the storage per instance, up to 4096GB. Can only be increased. e.g. --storage-size 20GB
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
      --user string               Username for the initial MongoDB user. Requires --password. Defaults to the cluster's current username
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The MongoDB version to upgrade to: one of the CanUpgradeTo versions of the current one. Downgrades are not supported. Asks for confirmation, as snapshots taken before the upgrade can no longer be restored in place
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mongo-v2 cluster update --location <location> --cluster-id <cluster-id> --cores 2 --ram 16GB --storage-size 20GB

# Scale to three instances, one address per instance
ionosctl dbaas mongo-v2 cluster update --location <location> --cluster-id <cluster-id> --instances 3 --cidr <cidr1>,<cidr2>,<cidr3>
```

