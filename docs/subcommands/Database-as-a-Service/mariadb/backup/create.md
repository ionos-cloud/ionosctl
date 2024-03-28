---
description: "Create an ad-hoc MariaDB Backup"
---

# DbaasMariadbBackupCreate

## Usage

```text
ionosctl dbaas mariadb backup create [flags]
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

For `create` command:

```text
[c]
```

## Description

Create an ad-hoc MariaDB Backup

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster (required)
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mariadb backup create --cluster-id CLUSTER_ID
```

