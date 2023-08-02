---
description: "List Cluster Backups from a Cluster"
---

# DbaasPostgresClusterBackupList

## Usage

```text
ionosctl dbaas postgres cluster backup list [flags]
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
                            Available columns: [BackupId ClusterId Active CreatedDate EarliestRecoveryTargetTime Version State] (default [BackupId,ClusterId,CreatedDate,EarliestRecoveryTargetTime,Active,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          When using text output, don't print headers
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas postgres backup list
```

