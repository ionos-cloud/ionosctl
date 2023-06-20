---
description: "Get a Cluster Backup"
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
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
  -i, --backup-id string   The unique ID of the Backup (required)
      --cols strings       Set of columns to be printed on output 
                           Available columns: [BackupId ClusterId Active CreatedDate EarliestRecoveryTargetTime Version State] (default [BackupId,ClusterId,CreatedDate,EarliestRecoveryTargetTime,Active,State])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --no-headers         When using text output, don't print headers
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
  -v, --verbose            Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas postgres backup get -i BACKUP_ID
```

