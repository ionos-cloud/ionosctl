---
description: Get a Cluster Backup
---

# DbaasPostgresBackupGet

## Usage

```text
ionosctl dbaas postgres backup get [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
```

For `backup` command:

```text
[b]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a PostgreSQL Backup by using its ID.

Required values to run command:

* Backup Id

## Options

```text
  -i, --backup-id string   The unique ID of the Backup (required)
      --no-headers         When using text output, don't print headers
```

## Examples

```text
ionosctl dbaas postgres backup get -i BACKUP_ID
```

