---
description: Delete a Snapshot
---

# Delete

## Usage

```text
ionosctl snapshot delete [flags]
```

## Description

Use this command to delete the specified Snapshot.

Required values to run command:

* Snapshot Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                Force command to execute without user input
  -h, --help                 help for delete
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
      --snapshot-id string   The unique Snapshot Id (required)
      --timeout int          Timeout option for a Snapshot to be deleted [seconds] (default 60)
      --wait                 Wait for Snapshot to be deleted
```

## Examples

```text
ionosctl snapshot delete --snapshot-id 8e0bc509-87ee-47f4-a382-302e4f7e103d --wait 
Warning: Are you sure you want to delete snapshot (y/N) ? 
y
RequestId: 6e029eb6-47e6-4dcd-a333-d620b49c01e5
Status: Command snapshot delete and request have been successfully executed
```

