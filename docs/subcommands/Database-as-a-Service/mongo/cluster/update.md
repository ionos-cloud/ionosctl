---
description: "Update a Mongo Cluster by ID"
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

Update a Mongo Cluster by ID

## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
      --backup-location string    The location where the cluster backups will be stored. If not set, the backup is stored in the nearest location of the cluster
      --biconnector string        The host and port where this new BI Connector is installed. The MongoDB Connector for Business Intelligence allows you to query a MongoDB database using SQL commands. Example: r1.m-abcdefgh1234.mongodb.de-fra.ionos.com:27015
      --biconnector-enabled       Enable or disable the biconnector. If left unset, no change will be made to the biconnector's status. To explicitly disable it, use --biconnector-enabled=false
      --cidr strings              The list of IPs and subnet for your cluster. All IPs must be in a /24 network. Note the following unavailable IP range: 10.233.114.0/24
  -i, --cluster-id string         The unique ID of the cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name Edition Type URL Instances Shards Health State MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                 The total number of cores for the Server, e.g. 4. (required and only settable for enterprise edition)
      --datacenter-id string      The datacenter to which your cluster will be connected. Must be in the same location as the cluster
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           The total number of instances of the cluster (one primary and n-1 secondaries). Minimum of 3 for business edition (default 1)
      --lan-id string             The numeric LAN ID with which you connect your cluster
      --maintenance-day string    Day for Maintenance. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: Saturday
      --maintenance-time string   Time for the Maintenance. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59
  -n, --name string               The name of your cluster
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
      --ram string                Custom RAM: multiples of 1024. e.g. --ram 1024 or --ram 1024MB or --ram 4GB (required and only settable for enterprise edition)
      --shards int32              The total number of shards in the sharded_cluster cluster. Setting this flag is only possible for enterprise clusters and infers a sharded_cluster type. Possible values: 2 - 32. (required for sharded_cluster enterprise clusters) (default 1)
      --storage-size string       Custom Storage: Greater performance for values greater than 100 GB. (required and only settable for enterprise edition)
      --storage-type string       Custom Storage Type. (only settable for enterprise edition) (default "\"SSD Standard\"")
  -v, --verbose                   Print step-by-step process when running command
  -w, --wait                      Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl dbaas mongo cluster update --cluster-id <cluster-id>
```

