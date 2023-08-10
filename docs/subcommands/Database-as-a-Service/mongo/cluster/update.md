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
[mongodb mdb m]
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
      --cidr strings              The list of IPs and subnet for your cluster. Note the following unavailable IP ranges: 10.233.114.0/24 (required)
  -i, --cluster-id string         The unique ID of the cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name URL Health State Instances MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string      The datacenter to which your cluster will be connected. Must be in the same location as the cluster (required)
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           The total number of instances in the cluster (one primary and n-1 secondaries)
      --lan-id string             The numeric LAN ID with which you connect your cluster (required)
      --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur (required)
      --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59 (required)
  -n, --name string               When using text output, don't print headers
      --no-headers                When using text output, don't print headers
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
      --template-id string        The unique ID of the template, which specifies the number of cores, storage size, and memory. You cannot downgrade to a smaller template or minor edition
  -v, --verbose                   Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo cluster update --cluster-id <cluster-id>
```

