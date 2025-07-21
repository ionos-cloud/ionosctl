---
description: "List Cluster Backups"
---

# DbaasPostgresBackupList

## Usage

```text
ionosctl dbaas postgres backup list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of PostgreSQL Cluster Backups.

## Options

```text
  -u, --api-url string   Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [BackupId ClusterId Active CreatedDate EarliestRecoveryTargetTime Version State] (default [BackupId,ClusterId,CreatedDate,EarliestRecoveryTargetTime,Active,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas postgres backup list
```

