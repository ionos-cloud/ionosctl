---
description: Get a Cluster Backup
---

# PgBackupGet

## Usage

```text
ionosctl pg backup get [flags]
```

## Aliases

For `pg` command:

```text
[postgres]
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
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
  -i, --backup-id string   The unique ID of the Backup (required)
      --cols strings       Set of columns to be printed on output 
                           Available columns: [BackupId ClusterId DisplayName Active CreatedDate EarliestRecoveryTargetTime Version] (default [BackupId,ClusterId,CreatedDate,EarliestRecoveryTargetTime,Active])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
  -v, --verbose            Print step-by-step process when running command
```

## Examples

```text
ionosctl pg backup get -i BACKUP_ID
```

