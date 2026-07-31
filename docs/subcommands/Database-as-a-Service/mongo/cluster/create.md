---
description: "Create DBaaS Mongo Replicaset or Sharded Clusters for your chosen edition"
---

# DbaasMongoClusterCreate

## Usage

```text
ionosctl dbaas mongo cluster create [flags]
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

For `create` command:

```text
[c]
```

## Description

Create a MongoDB cluster. The required flags depend on the --edition you pick, because each edition bounds the topology and how sizing is expressed:

  playground  - Fixed single-instance sandbox. Type is forced to replicaset, --instances to 1, and the Playground template (1 core / 2 GB RAM / 50 GB storage) is applied automatically. No snapshots/restore. Only the always-required flags below are needed.
  business    - Replicaset only. Sizing comes from a --template (bundle of cores/RAM/storage); defaults to the XS template if omitted. You must pass --instances (1 or 3).
  enterprise  - Replicaset or sharded-cluster with EXPLICIT sizing. Templates are forbidden here; instead set --cores, --ram, --storage-size, --storage-type (all default to the smallest valid value). Enterprise unlocks sharding and point-in-time restore.

Always required (every edition): --name, --datacenter-id, --lan-id, --cidr.

Choosing type for enterprise: pass --type explicitly, or let it be inferred — setting --shards implies 'sharded-cluster', otherwise 'replicaset'. A sharded-cluster requires --shards (2-32); a replicaset requires --instances.

The cluster is provisioned into the same location as its datacenter (--location is inferred from --datacenter-id when omitted). All CIDRs in --cidr must be /24 and there must be one per instance.

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --backup-location string    Object Storage region where snapshots (backups) are stored, e.g. de, eu-south-2, eu-central-3. For extra safety pick a region different from the cluster. Defaults to the region nearest the cluster (eu-central-2 where S3 is unavailable)
      --biconnector string        Enable the MongoDB BI Connector at this host:port so BI tools can query the cluster over SQL. Example: r1.m-abcdefgh1234.mongodb.de-fra.ionos.com:27015. Setting this turns the BI Connector on
      --biconnector-enabled       Whether the BI Connector configured via --biconnector is enabled. Pass --biconnector-enabled=false to keep it configured but off (default true)
      --cidr strings              Comma-separated private IPs (with /24 subnet) assigned to the cluster nodes on the LAN - one per instance. Each must be a /24, e.g. 192.168.1.100/24,192.168.1.101/24. Unavailable range: 10.233.114.0/24 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name Edition Type URL Instances Shards Health State MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int32               CPU cores PER instance, e.g. 4. Minimum 1. Enterprise only (playground/business derive this from their template) (default 1)
      --datacenter-id string      Virtual Datacenter the cluster attaches to. The cluster is reachable privately from this VDC and inherits its location. Must be in the same location as the cluster (required)
  -D, --depth int                 Level of detail for response objects (default 1)
  -e, --edition string            The tier that bounds topology and sizing: playground (fixed 1-instance sandbox, no backups), business (replicaset sized by a template), enterprise (replicaset or sharded-cluster with explicit cores/RAM/storage and point-in-time restore). Can be inferred from --template. Can be one of: playground, business, enterprise (required)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           Number of nodes in the replicaset: one primary and n-1 secondaries. Use an odd count so a majority can elect a primary. Valid values: 1, 3, 5, 7 (playground is fixed at 1; business allows 1 or 3). Required for business and for enterprise replicaset (default 1)
      --lan-id string             The private LAN (within --datacenter-id) the cluster connects to, given as its numeric LAN ID (required)
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           The physical location (region) where the cluster is created, e.g. de/fra. Immutable. Defaults to the location of the connected --datacenter-id; the datacenter and cluster must share the same location
      --maintenance-day string    Day of the week for the weekly 4-hour maintenance window (e.g. Saturday). Paired with --maintenance-time. Defaults to a random weekday (Mon-Fri) (default "Random (Mon-Fri 10:00-16:00)")
      --maintenance-time string   Start time (UTC, HH:MM:SS) of the weekly 4-hour maintenance window during which IONOS may apply updates/patches, e.g. 16:30:59. Paired with --maintenance-day. Defaults to a random time between 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
  -n, --name string               The human-readable display name of your cluster (shown in listings and the DCD) (required)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                RAM PER instance. Must be a multiple of 1024 MB (1 GB); minimum 2 GB. Accepts a unit, e.g. --ram 2048, --ram 2048MB or --ram 4GB. Enterprise only (default "2GB")
      --shards int32              Number of shards the data is partitioned across (each shard is itself a replicaset). Valid values: 2-32. Enterprise only; setting it infers --type sharded-cluster and is required for a sharded-cluster (required) (default 1)
      --storage-size string       Disk size PER instance. Accepts a unit, e.g. --storage-size 100GB. Minimum 2 GB; noticeably better performance above 100 GB. Enterprise only (default "5GB")
      --storage-type string       Disk type backing each instance: HDD, 'SSD Standard' or 'SSD Premium' (fastest). Enterprise only (default "\"SSD Standard\"")
      --template templates list   Name (e.g. XS, S, L, 4XL) or ID of a Mongo Template - a predefined bundle of cores, RAM and storage size (see templates list). Used for playground/business editions; forbidden for enterprise (use --cores/--ram/--storage-size instead). Business defaults to the XS template, playground to the Playground template. Setting a template can also infer --edition.
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
      --type string               replicaset (one primary + n-1 identical secondaries) or sharded-cluster (data partitioned across shards; enterprise only). Required for enterprise unless inferred: setting --shards implies sharded-cluster, otherwise replicaset. Can be one of: replicaset, sharded-cluster (default "replicaset")
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version api-versions      The MongoDB major version to run, e.g. 6.0 or 7.0. Patch versions are managed automatically. Use api-versions/completion to see supported versions (required) (default "7.0")
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mongo cluster create --edition playground --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR 

ionosctl dbaas mongo cluster create --edition business --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --instances INSTANCES 

ionosctl dbaas mongo cluster create --edition enterprise --instances INSTANCES [--shards SHARDS] --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR 

ionosctl dbaas mongo cluster create --edition enterprise --type replicaset --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --instances INSTANCES 

ionosctl dbaas mongo cluster create --edition enterprise --type sharded-cluster --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --shards SHARDS 
```

