---
description: "Use a MariaDB Backup to restore the cluster to its previous state"
---

# DbaasMariadbBackupRestore

## Usage

```text
ionosctl dbaas mariadb backup restore [flags]
```

## Aliases

For `mariadb` command:

```text
[maria mar]
```

For `backup` command:

```text
[b]
```

For `restore` command:

```text
[r rs]
```

## Description

Use a MariaDB Backup to restore the cluster to its previous state

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -M, --max-results int32   The maximum number of elements to return
  -n, --name string         Response filter to list only the MariaDB Clusters that contain the specified name in the DisplayName field. The value is case insensitive
      --no-headers          Don't print table headers when table output is used
      --offset int32        Skip a certain number of results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mariadb backup restore --cluster-id CLUSTER_ID --backup-id BACKUP_ID
```

