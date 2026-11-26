---
description: "Create a MongoDB Cluster"
---

# DbaasMongoV2ClusterCreate

## Usage

```text
ionosctl dbaas mongo-v2 cluster create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create a new MongoDB Cluster with one primary and n-1 secondary instances. The --edition decides how the instances are sized:
  1. ENTERPRISE: pass --cores, --ram and --storage-size (defaults: 1 core, 8GB RAM, 10GB storage). Optionally enable sharding with --shards together with more than one replica per shard via --instances (sharding is permanent: shards can be added later, but never removed, and sharding cannot be disabled).
  2. BUSINESS: pass --template-id, the ID of an instance template that bundles cores, RAM and storage (see `ionosctl dbaas mongo-v2 template list`). The template cannot be changed later.

Both editions require the connection and credential flags (--datacenter-id, --lan-id, --cidr, --user, --password, --database) plus --location. --cidr takes one IP address with netmask per instance; with sharding it takes exactly three, one per config server, whatever the instance and shard counts. The snapshot location (--backup-location) cannot be changed later either.

To create the cluster from an existing snapshot instead of empty, additionally pass --snapshot-id. The cluster version is taken from the snapshot (so --version is not needed; if given, it must match the snapshot's version). The snapshot is applied in full, up to the end of its recovery window, unless --recovery-time picks an earlier point within that window. Each instance needs at least the snapshot's RequiredStorageSize of storage: that value is in MiB, while --storage-size and template storage are in GB (1GB = 1024MiB, round up).

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mongodbv2' and env var 'IONOS_API_URL' (default "https://mongodb.%s.ionos.com")
  -B, --backup-location string    The Object Storage location where snapshots will be stored. Cannot be changed later. For added data safety, use a different location than the cluster (default "eu-central-3")
      --cidr strings              The IP addresses with netmask of the cluster instances: one per instance, or exactly three with sharding. e.g.: 192.168.1.100/24,192.168.1.101/24,192.168.1.102/24. The range 10.208.0.0/12 is not available (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name Edition Version Instances State ConnectionUri Cores Ram StorageSize TemplateId ShardingEnabled Shards Description SnapshotLocation RetentionDays MaintenanceDay MaintenanceTime LogsEnabled MetricsEnabled DatacenterId LanId Cidr Username Database StatusMessage]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int32               ENTERPRISE only: the number of CPU cores per instance, between 1 and 62 (default 1)
      --database string           The name of the initial database, owned by the initial user: 2-63 letters, digits or underscores (required)
      --datacenter-id string      The unique ID of the Datacenter to connect to your cluster. Must be in the same location as the cluster (required)
  -D, --depth int                 Level of detail for response objects (default 1)
      --description string        Human-readable description for the cluster
      --edition string            The edition of the cluster, which cannot be changed later. ENTERPRISE is sized via --cores, --ram and --storage-size and supports sharding; BUSINESS is sized by an instance template via --template-id. Can be one of: ENTERPRISE, BUSINESS (required)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           The total number of instances in the cluster (one primary and n-1 secondary), between 1 and 9. With sharding, this is the number of replicas per shard, and must be more than 1 (default 1)
      --lan-id string             The unique ID of the LAN to connect your cluster to (required)
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: gb/bhx, us/ewr, gb/lhr, fr/par, es/vit, us/mci, de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --logs-enabled              Forward this cluster's logs to the IONOS Logging Service. Has no effect unless that service is activated on your contract, which incurs additional costs
      --maintenance-day string    Day of the week for the MaintenanceWindow. Defaults to a random day during Mon-Fri. Can be one of: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday (default "Random (Mon-Fri 10:00-16:00)")
      --maintenance-time string   Time for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59. Defaults to a random time during 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
      --metrics-enabled           Forward this cluster's metrics to the IONOS Monitoring Service. Has no effect unless that service is activated on your contract, which incurs additional costs
  -n, --name string               The friendly name of your cluster (default "UnnamedCluster")
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --password string           The password of the initial MongoDB user: 8-256 characters (required)
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                ENTERPRISE only: the amount of memory per instance, between 6GB and 240GB. e.g. --ram 8, --ram 8GB (default "8GB")
  -R, --recovery-time string      Advanced: with --snapshot-id, restore to a specific point in time WITHIN the snapshot's recovery window (PITR): 'now', a date, a date-time, or an RFC3339 timestamp (no timezone = UTC). Defaults to the end of the window, i.e. the snapshot applied in full
      --retention-days int32      How many days snapshots are retained before being automatically deleted. One of: 3, 4, 5, 6, 7, 15, 30, 60, 90, 120, 180, 360 (default 7)
      --shards int32              ENTERPRISE only: enable sharding and split the data across this many shards, between 2 and 5. --instances then sets the replicas per shard and must be more than 1, and --cidr takes exactly three addresses, one per config server. Sharding is permanent: shards can be added later but never removed, and it cannot be disabled
      --snapshot-id string        Create the cluster from this snapshot instead of empty. The snapshot must come from a cluster in the same location, and each instance needs at least its RequiredStorageSize (MiB) of storage. The connection/credential flags and --location are still required; the cluster version is taken from the snapshot
      --storage-size string       ENTERPRISE only: the size of the storage per instance, between 10GB and 4096GB. e.g. --storage-size 10 or --storage-size 10GB (default "10GB")
      --template-id string        BUSINESS only, and required there: the ID of the instance template that sets the cores, RAM and storage of each instance. Cannot be changed later. See 'ionosctl dbaas mongo-v2 template list'
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
      --user string               The username of the initial MongoDB user: 2-16 letters, digits or underscores (required)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The MongoDB version of your cluster. Defaults to the snapshot's version when --snapshot-id is set (default "7.0")
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create an empty ENTERPRISE cluster
ionosctl dbaas mongo-v2 cluster create --location <location> --edition ENTERPRISE --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --database <database>

# Create a sharded ENTERPRISE cluster: 2 shards of 3 replicas each, plus 3 addresses for the config servers
ionosctl dbaas mongo-v2 cluster create --location <location> --edition ENTERPRISE --shards 2 --instances 3 --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr1>,<cidr2>,<cidr3> --user <username> --password <password> --database <database>

# Create a BUSINESS cluster with one instance, sized by an instance template
ionosctl dbaas mongo-v2 cluster create --location <location> --edition BUSINESS --template-id <template-id> --instances 1 --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --database <database>

# Create a cluster from an existing snapshot, restored to a point inside its recovery window (version is taken from the snapshot)
ionosctl dbaas mongo-v2 cluster create --location <location> --edition ENTERPRISE --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --database <database> --snapshot-id <snapshot-id> --recovery-time 2025-01-02T15:00:00Z
```

