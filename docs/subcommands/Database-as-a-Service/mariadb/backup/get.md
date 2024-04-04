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
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --backup-id string   The ID of the Backup to be retrieved
      --cols strings       Set of columns to be printed on output 
                           Available columns: [BackupId ClusterId Size Items] (default [BackupId,ClusterId,Size,Items])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -v, --verbose            Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mariadb backup get --backup-id BACKUP_ID
```

