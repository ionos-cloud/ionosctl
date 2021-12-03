---
description: List Cluster Backups from a Cluster
---

# PgClusterBackupList

## Usage

```text
ionosctl pg cluster backup list [flags]
```

## Aliases

For `cluster` command:

```text
[c]
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

Use this command to retrieve a list of PostgreSQL Cluster Backups from a specific Cluster.

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [BackupId ClusterId DisplayName Active CreatedDate EarliestRecoveryTargetTime Version] (default [BackupId,ClusterId,DisplayName,EarliestRecoveryTargetTime,Active])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl pg backup list
```

