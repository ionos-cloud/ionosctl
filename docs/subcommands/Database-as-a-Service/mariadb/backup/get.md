---
description: "Get a MariaDB Backup"
---

# DbaasMariadbBackupGet

## Usage

```text
ionosctl dbaas mariadb backup get [flags]
```

## Aliases

For `mariadb` command:

```text
[maria mar ma]
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

Get a MariaDB Backup

## Options

```text
  -u, --api-url string     Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mariadb' and env var 'IONOS_API_URL' (default "https://mariadb.%s.ionos.com")
      --backup-id string   The ID of the Backup to be retrieved (required)
      --cols strings       Set of columns to be printed on output 
                           Available columns: [BackupId ClusterId Size Items] (default [BackupId,ClusterId,Size,Items])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --limit int          Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string    Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci (default "de/txl")
      --no-headers         Don't print table headers when table output is used
      --offset int         Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string      Desired output format [text|json|api-json] (default "text")
      --query string       JMESPath query string to filter the output
  -q, --quiet              Quiet output
  -v, --verbose count      Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas mariadb backup get --backup-id BACKUP_ID
```

