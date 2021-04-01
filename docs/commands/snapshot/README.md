---
description: Snapshot Operations
---

# Snapshot

## Usage

```text
ionosctl snapshot [command]
```

## Description

The sub-command of `ionosctl snapshot` allows you to see information about snapshots available.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for snapshot
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl snapshot create](create.md) | Create a Snapshot |
| [ionosctl snapshot get](get.md) | Get a Snapshot |
| [ionosctl snapshot list](list.md) | List Snapshots |

