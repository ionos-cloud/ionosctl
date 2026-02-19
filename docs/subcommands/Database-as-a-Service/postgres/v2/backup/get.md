---
description: "Get a PostgreSQL Backup"
---

# DbaasPostgresV2BackupGet

## Usage

```text
ionosctl dbaas postgres-v2 backup get [flags]
```

## Aliases

For `postgres-v2` command:

```text
[pg-v2]
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
  -u, --api-url string     Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'todo' and env var 'IONOS_API_URL' (default "https://postgresql.%s.ionos.com")
  -i, --backup-id string   The unique ID of the Backup (required)
      --cols strings       Set of columns to be printed on output 
                           Available columns: [BackupId ClusterId PostgresClusterVersion Location IsActive EarliestRecoveryTargetTime LatestRecoveryTargetTime State CreatedDate] (default [BackupId,ClusterId,PostgresClusterVersion,Location,IsActive,EarliestRecoveryTargetTime,LatestRecoveryTargetTime])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int          Level of detail for response objects (default 1)
  -F, --filters strings    Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --limit int          Maximum number of items to return per request (default 50)
  -l, --location string    Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, gb/bhx, us/las, us/mci, us/ewr (default "de/txl")
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
ionosctl dbaas postgres backup get --backup-id <backup-id>
```

