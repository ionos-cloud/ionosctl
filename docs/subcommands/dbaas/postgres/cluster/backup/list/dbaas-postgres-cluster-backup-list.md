---
description: List Cluster Backups from a Cluster
---

# DbaasPostgresClusterBackupList

## Usage

```text
ionosctl dbaas postgres cluster backup list [flags]
```

## Aliases

For `cluster` command:

```text
[c]
```

For `backup` command:

```text
[b]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of PostgreSQL Cluster Backups from a specific Cluster.

## Options

```text
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [BackupId ClusterId Active CreatedDate EarliestRecoveryTargetTime Version State] (default [BackupId,ClusterId,CreatedDate,EarliestRecoveryTargetTime,Active,State])
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl dbaas postgres backup list
```

