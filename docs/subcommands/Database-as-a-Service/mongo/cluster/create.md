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
      --cidr strings              The list of IPs and subnet for your cluster. All IPs must be in a /24 network. Note the following unavailable IP range: 10.233.114.0/24 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name URL State Instances MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string      The datacenter to which your cluster will be connected. Must be in the same location as the cluster (required)
  -e, --edition string            Cluster Edition. Can be one of: playground, business, enterprise (required)
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           The total number of instances in the cluster (one primary and n-1 secondaries). (required for non-playground clusters) (default 1)
      --lan-id string             The numeric LAN ID with which you connect your cluster (required)
  -l, --location string           The physical location where the cluster will be created. (defaults to the location of the connected datacenter)
      --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: Saturday (required)
      --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59 (required)
  -n, --name string               The name of your cluster
      --no-headers                When using text output, don't print headers
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
      --template-id string        The unique ID of the template, which specifies the number of cores, storage size, and memory. Must not be provided for enterprise clusters
  -t, --timeout int               Timeout option for Request [seconds] (default 60)
      --type string               Cluster Type. Required for enterprise clusters. Can be one of: replicaset, sharded-cluster (default "replicaset")
  -v, --verbose                   Print step-by-step process when running command
      --version string            The MongoDB version of your cluster (default "6.0")
  -w, --wait-for-request          Wait for the Request to be executed
```

## Examples

```text
ionosctl dbaas mongo cluster create --edition playground --name NAME --maintenance-day MAINTENANCE_DAY --maintenance-time MAINTENANCE_TIME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR 

ionosctl dbaas mongo cluster create --edition business --name NAME --maintenance-day MAINTENANCE_DAY --maintenance-time MAINTENANCE_TIME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --template-id TEMPLATE_ID --instances INSTANCES 

ionosctl dbaas mongo cluster create --edition enterprise (--instances INSTANCES | --shards SHARDS) --name NAME --maintenance-day MAINTENANCE_DAY --maintenance-time MAINTENANCE_TIME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --cores CORES --storage-type STORAGE_TYPE --storage-size STORAGE_SIZE --ram RAM 

ionosctl dbaas mongo cluster create --edition enterprise --type replicaset --name NAME --maintenance-day MAINTENANCE_DAY --maintenance-time MAINTENANCE_TIME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --cores CORES --storage-type STORAGE_TYPE --storage-size STORAGE_SIZE --ram RAM --instances INSTANCES 

ionosctl dbaas mongo cluster create --edition enterprise --type sharded-cluster --name NAME --maintenance-day MAINTENANCE_DAY --maintenance-time MAINTENANCE_TIME --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --cores CORES --storage-type STORAGE_TYPE --storage-size STORAGE_SIZE --ram RAM --shards SHARDS 
```

