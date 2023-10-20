---
description: "Restore a PostgreSQL Cluster"
---

# DbaasPostgresClusterRestore

## Usage

```text
ionosctl dbaas postgres cluster restore [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
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

Use this command to trigger an in-place restore of the specified PostgreSQL Cluster.

Required values to run command:

* Cluster Id
* Backup Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --backup-id string       The unique ID of the backup you want to restore (required)
  -i, --cluster-id string      The unique ID of the Cluster (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ClusterId DisplayName Location State PostgresVersion Instances Ram Cores StorageSize StorageType DatacenterId LanId Cidr MaintenanceWindow SynchronizationMode BackupLocation] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,Cidr,Instances,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -R, --recovery-time string   If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely
  -t, --timeout int            Timeout option for Cluster to be in AVAILABLE state[seconds] (default 1200)
  -v, --verbose                Print step-by-step process when running command
  -W, --wait-for-state         Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl dbaas postgres cluster restore -i CLUSTER_ID --backup-id BACKUP_ID
```

