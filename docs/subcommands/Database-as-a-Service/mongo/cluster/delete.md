---
description: "Delete a Mongo Cluster by ID"
---

# DbaasMongoClusterDelete

## Usage

```text
ionosctl dbaas mongo cluster delete [flags]
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

For `delete` command:

```text
[del d]
```

## Description

Delete a Mongo Cluster by ID

## Options

```text
  -a, --all                 Delete all mongo clusters
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId Name Edition Type URL Instances Shards Health State MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --name                When deleting all clusters, filter the clusters by a name
      --no-headers          When using text output, don't print headers
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo cluster delete --cluster-id <cluster-id>
ionosctl db m c d --all
ionosctl db m c d --all --name <name>
```

