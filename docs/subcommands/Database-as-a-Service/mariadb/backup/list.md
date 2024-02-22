---
description: "List MariaDB Backups"
---

# DbaasMariadbBackupList

## Usage

```text
ionosctl dbaas mariadb backup list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of MariaDB Backups provisioned under your account.

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
ionosctl dbaas mariadb backup list
```

