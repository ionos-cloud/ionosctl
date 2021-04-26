---
description: Snapshot Operations
---

# Snapshot

## Usage

```text
ionosctl snapshot [command]
```

## Description

The sub-commands of `ionosctl snapshot` allow you to see information, to create, update, delete Snapshots.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for snapshot
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl snapshot add-label](add-label.md) | Add a Label on a Snapshot |
| [ionosctl snapshot create](create.md) | Create a Snapshot of a Volume within the Virtual Data Center. |
| [ionosctl snapshot delete](delete.md) | Delete a Snapshot |
| [ionosctl snapshot get](get.md) | Get a Snapshot |
| [ionosctl snapshot get-label](get-label.md) | Get a Label from a Snapshot |
| [ionosctl snapshot list](list.md) | List Snapshots |
| [ionosctl snapshot list-labels](list-labels.md) | List Labels from a Snapshot |
| [ionosctl snapshot remove-label](remove-label.md) | Remove a Label from a Snapshot |
| [ionosctl snapshot restore](restore.md) | Restore a Snapshot onto a Volume |
| [ionosctl snapshot update](update.md) | Update a Snapshot. |

