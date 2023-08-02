---
description: "Restore a Mongo Cluster by ID, using a snapshot"
---

# DbaasMongoClusterRestore

## Usage

```text
ionosctl dbaas mongo cluster restore [flags]
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

For `restore` command:

```text
[r]
```

## Description

Restore a Mongo Cluster by ID, using a snapshot

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string    The unique ID of the cluster (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ClusterId Name URL State Instances MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           When using text output, don't print headers
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
      --snapshot-id string   The unique ID of the snapshot you want to restore. (required)
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo cluster restore --cluster-id <cluster-id> --snapshot-id <snapshot-id>
```

