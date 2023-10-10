---
description: "Restore a Mongo Cluster using a snapshot"
---

# DbaasMongoSnapshotRestore

## Usage

```text
ionosctl dbaas mongo snapshot restore [flags]
```

## Aliases

For `mongo` command:

```text
[mongodb mdb m]
```

For `snapshot` command:

```text
[snap backup snapshots backups]
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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string    The unique ID of the cluster (required)
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
      --snapshot-id string   The unique ID of the snapshot you want to restore. (required)
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo cluster restore --cluster-id <cluster-id> --snapshot-id <snapshot-id>
```

