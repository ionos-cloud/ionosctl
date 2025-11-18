---
description: "Restore a Mongo Cluster using a snapshot"
---

# DbaasMongoClusterRestore

## Usage

```text
ionosctl dbaas mongo cluster restore [flags]
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

For `restore` command:

```text
[r]
```

## Description

This command restores a cluster via its snapshot. A cluster can have multiple snapshots. A snapshot is added during the following cases:
When a cluster is created, known as initial sync which usually happens in less than 24 hours.
After a restore.
A snapshot is a copy of the data in the cluster at a certain time. Every 24 hours, a base snapshot is taken, and every Sunday, a full snapshot is taken. Snapshots are retained for the last seven days; hence, recovery is possible for up to a week from the current date.
You can restore from any snapshot as long as it was created with the same or older MongoDB patch version.
Snapshots are stored in an IONOS S3 Object Storage bucket in the same region as your database. Databases in regions where IONOS S3 Object Storage is not available is backed up to eu-central-2.

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string    The unique ID of the cluster (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ClusterId Name Edition Type URL Instances Shards Health State MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
      --snapshot-id string   The unique ID of the snapshot you want to restore. (required)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas mongo cluster restore --cluster-id <cluster-id> --snapshot-id <snapshot-id>
```

