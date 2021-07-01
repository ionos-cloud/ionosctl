---
description: List Snapshots
---

# SnapshotList

## Usage

```text
ionosctl snapshot list [flags]
```

## Aliases

For `snapshot` command:

```text
[ss snap]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Snapshots.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [SnapshotId Name LicenceType Size State] (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl snapshot list
```

