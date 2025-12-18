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
  -u, --api-url string     Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
  -i, --backup-id string   The unique ID of the Backup (required)
      --cols strings       Set of columns to be printed on output 
                           Available columns: [BackupId ClusterId Active CreatedDate EarliestRecoveryTargetTime Version State] (default [BackupId,ClusterId,CreatedDate,EarliestRecoveryTargetTime,Active,State])
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
  -v, --verbose count      Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas postgres backup get -i BACKUP_ID
```

