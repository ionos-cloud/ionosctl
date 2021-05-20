---
description: Delete a Snapshot
---

# SnapshotDelete

## Usage

```text
ionosctl snapshot delete [flags]
```

## Aliases

For `snapshot` command:
```text
[snap]
```

## Description

Use this command to delete the specified Snapshot.

Required values to run command:

* Snapshot Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -F, --format strings       Collection of fields to be printed on output (default [SnapshotId,Name,LicenceType,Size,State])
  -h, --help                 help for delete
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
      --snapshot-id string   The unique Snapshot Id (required)
  -t, --timeout int          Timeout option for Request for Snapshot deletion [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Snapshot deletion to be executed
```

## Examples

```text
ionosctl snapshot delete --snapshot-id 8e0bc509-87ee-47f4-a382-302e4f7e103d --wait-for-request 
Warning: Are you sure you want to delete snapshot (y/N) ? 
y
RequestId: 6e029eb6-47e6-4dcd-a333-d620b49c01e5
Status: Command snapshot delete & wait have been successfully executed
```

