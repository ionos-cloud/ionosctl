---
description: "Get a backup"
---

# DbaasPostgresBackupGet

## Usage

```text
ionosctl dbaas postgres backup get [flags]
```

## Aliases

For `postgres` command:

```text
[pg pgsql postgresql psql]
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

Retrieve details of a single backup by its ID, including the cluster it belongs to, its PostgreSQL version, and the start of its recovery window (EarliestRecoveryTargetTime) which bounds the --recovery-time you can use when restoring.

Required values to run command:

* Backup Id

## Options

```text
  -u, --api-url string     Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
  -i, --backup-id string   ID of the backup to retrieve. See 'dbaas postgres backup list' (required)
      --cols strings       Set of columns to be printed on output 
                           Available columns: [BackupId ClusterId Active CreatedDate EarliestRecoveryTargetTime Version State]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int          Level of detail for response objects (default 1)
  -F, --filters strings    Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --limit int          Maximum number of items to return per request (default 50)
      --no-headers         Don't print table headers when table output is used
      --offset int         Number of items to skip before starting to collect the results
      --order-by string    Property to order the results by
  -o, --output string      Desired output format [text|json|api-json] (default "text")
      --query string       JMESPath query string to filter the output
  -q, --quiet              Quiet output
  -t, --timeout int        Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count      Increase verbosity level [-v, -vv, -vvv]
  -w, --wait               Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas postgres backup get -i BACKUP_ID
```

