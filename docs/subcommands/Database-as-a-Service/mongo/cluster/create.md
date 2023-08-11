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
[mongodb mdb m]
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
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
      --backup-location string    The location where the cluster backups will be stored. If not set, the backup is stored in the nearest location of the cluster
      --biconnector string        The host and port where this new BI Connector is installed. The MongoDB Connector for Business Intelligence allows you to query a MongoDB database using SQL commands. Example: r1.m-abcdefgh1234.mongodb.de-fra.ionos.com:27015
      --cidr strings              The list of IPs and subnet for your cluster. All IPs must be in a /24 network. Note the following unavailable IP range: 10.233.114.0/24 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name Edition Type URL Instances Shards Health State MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                 The total number of cores for the Server, e.g. 4. (required and only settable for enterprise edition) (required)
      --datacenter-id string      The datacenter to which your cluster will be connected. Must be in the same location as the cluster (required)
  -e, --edition string            Cluster Edition. Can be one of: playground, business, enterprise (required)
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           The total number of instances of the cluster (one primary and n-1 secondaries). Minimum of 3 for business edition (default 1)
      --lan-id string             The numeric LAN ID with which you connect your cluster (required)
  -l, --location string           The physical location where the cluster will be created. (defaults to the location of the connected datacenter)
      --maintenance-day string    Day for Maintenance. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: Saturday (required)
      --maintenance-time string   Time for the Maintenance. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59 (required)
  -n, --name string               The name of your cluster (required)
      --no-headers                When using text output, don't print headers
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
      --ram string                Custom RAM: multiples of 1024. e.g. --ram 1024 or --ram 1024MB or --ram 4GB (required and only settable for enterprise edition) (required)
      --shards int32              The total number of shards in the sharded_cluster cluster. Setting this flag is only possible for enterprise clusters and infers a sharded_cluster type. Possible values: 2 - 32. (required for sharded_cluster enterprise clusters) (required) (default 1)
      --storage-size string       Custom Storage: Greater performance for values greater than 100 GB. (required and only settable for enterprise edition) (required)
      --storage-type string       Custom Storage Type. (required and only settable for enterprise edition). Can be one of: HDD, SSD, "SSD Premium" (required)
      --template string           The ID of a Mongo Template, or a word contained in the name of one. Templates specify the number of cores, storage size, and memory. (Required only for business edition) (required)
  -t, --timeout int               Timeout option for Request [seconds] (default 60)
      --type string               Cluster Type. Required for enterprise clusters. Not required (inferred) if using --shards or --instances. Can be one of: replicaset, sharded-cluster (default "replicaset")
  -v, --verbose                   Print step-by-step process when running command
      --version string            The MongoDB version of your cluster (required) (default "6.0")
  -w, --wait-for-request          Wait for the Request to be executed
```

## Examples

```text
ionosctl dbaas mongo cluster create --edition playground --name NAME --maintenance-day MAINTENANCE_DAY --maintenance-time MAINTENANCE_TIME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR 

ionosctl dbaas mongo cluster create --edition business --name NAME --maintenance-day MAINTENANCE_DAY --maintenance-time MAINTENANCE_TIME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --template TEMPLATE --instances INSTANCES 

ionosctl dbaas mongo cluster create --edition enterprise (--instances INSTANCES | --shards SHARDS) --name NAME --maintenance-day MAINTENANCE_DAY --maintenance-time MAINTENANCE_TIME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --cores CORES --storage-type STORAGE_TYPE --storage-size STORAGE_SIZE --ram RAM 

ionosctl dbaas mongo cluster create --edition enterprise --type replicaset --name NAME --maintenance-day MAINTENANCE_DAY --maintenance-time MAINTENANCE_TIME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --cores CORES --storage-type STORAGE_TYPE --storage-size STORAGE_SIZE --ram RAM --instances INSTANCES 

ionosctl dbaas mongo cluster create --edition enterprise --type sharded-cluster --name NAME --maintenance-day MAINTENANCE_DAY --maintenance-time MAINTENANCE_TIME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --cores CORES --storage-type STORAGE_TYPE --storage-size STORAGE_SIZE --ram RAM --shards SHARDS 
```

