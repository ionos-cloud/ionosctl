---
description: Create Mongo Clusters
---

# MongoClusterCreate

## Usage

```text
mongo cluster create [flags]
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

Create Mongo Clusters

## Options

```text
      --cidr strings              The list of IPs and subnet for your cluster. Note the following unavailable IP ranges: 10.233.64.0/18 10.233.0.0/18 10.233.114.0/24 (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name URL State Instances MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId]
      --datacenter-id string      The datacenter to which your cluster will be connected. Must be in the same location as the cluster (required)
  -h, --help                      help for create
      --instances int32           The total number of instances in the cluster (one primary and n-1 secondaries)
      --lan-id string             The numeric LAN ID with which you connect your cluster (required)
  -l, --location string           The physical location where the cluster will be created. Defaults to the connection's datacenter location
      --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur (required)
      --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59 (required)
      --mongo-version string      The MongoDB version of your cluster (default "5.0")
  -n, --name string               The name of your cluster
      --no-headers                When using text output, don't print headers
      --template-id string        The unique ID of the template, which specifies the number of cores, storage size, and memory
  -t, --timeout int               Timeout option for Request [seconds] (default 60)
  -w, --wait-for-request          Wait for the Request to be executed
```

