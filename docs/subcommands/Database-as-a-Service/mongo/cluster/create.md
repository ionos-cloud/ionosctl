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

Create DBaaS Mongo Replicaset or Sharded Clusters for your chosen edition

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --backup-location string    The location where the cluster backups will be stored. If not set, the backup is stored in the nearest location of the cluster
      --biconnector string        BI Connector host & port. The MongoDB Connector for Business Intelligence allows you to query a MongoDB database using SQL commands. Example: r1.m-abcdefgh1234.mongodb.de-fra.ionos.com:27015
      --biconnector-enabled       Enable or disable the biconnector. To disable it, use --biconnector-enabled=false (default true)
      --cidr strings              The list of IPs and subnet for your cluster. All IPs must be in a /24 network (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name Edition Type URL Instances Shards Health State MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int32               The total number of cores for the Server, e.g. 4. (required and only settable for enterprise edition) (default 1)
      --datacenter-id string      The datacenter to which your cluster will be connected. Must be in the same location as the cluster (required)
  -e, --edition string            Cluster Edition. Can be one of: playground, business, enterprise (required)
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           The total number of instances of the cluster (one primary and n-1 secondaries). Minimum of 3 for enterprise edition (default 1)
      --lan-id string             The numeric LAN ID with which you connect your cluster (required)
      --limit int                 Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string           The physical location where the cluster will be created. (defaults to the location of the connected datacenter)
      --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. Defaults to a random day during Mon-Fri, during the hours 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
      --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59. Defaults to a random day during Mon-Fri, during the hours 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
  -n, --name string               The name of your cluster (required)
      --no-headers                Don't print table headers when table output is used
      --offset int                Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram string                Custom RAM: multiples of 1024. e.g. --ram 1024 or --ram 1024MB or --ram 4GB (required and only settable for enterprise edition) (default "2GB")
      --shards int32              The total number of shards in the sharded_cluster cluster. Setting this flag is only possible for enterprise clusters and infers a sharded_cluster type. Possible values: 2 - 32. (required for sharded_cluster enterprise clusters) (required) (default 1)
      --storage-size string       Custom Storage: Minimum of 5GB, Greater performance for values greater than 100 GB. (only settable for enterprise edition) (default "5GB")
      --storage-type string       Custom Storage Type. (only settable for enterprise edition) (default "\"SSD Standard\"")
      --template string           The ID of a Mongo Template, or a word contained in the name of one. Templates specify the number of cores, storage size, and memory. Business editions default to XS template. Playground editions default to playground template.
      --type string               Cluster Type. Required for enterprise clusters. Not required (inferred) if using --shards or --instances. Can be one of: replicaset, sharded-cluster (default "replicaset")
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version string            The MongoDB version of your cluster (required) (default "7.0")
```

## Examples

```text
ionosctl dbaas mongo cluster create --edition playground --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR 

ionosctl dbaas mongo cluster create --edition business --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --instances INSTANCES 

ionosctl dbaas mongo cluster create --edition enterprise --instances INSTANCES [--shards SHARDS] --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR 

ionosctl dbaas mongo cluster create --edition enterprise --type replicaset --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --instances INSTANCES 

ionosctl dbaas mongo cluster create --edition enterprise --type sharded-cluster --name NAME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --shards SHARDS 
```

