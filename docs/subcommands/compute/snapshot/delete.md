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
[ss snap]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete the specified Snapshot.

Required values to run command:

* Snapshot Id

## Options

```text
  -a, --all                  Delete all the Snapshots.
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -i, --snapshot-id string   The unique Snapshot Id (required)
  -t, --timeout int          Timeout option for Request for Snapshot deletion [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Snapshot deletion to be executed
```

## Examples

```text
ionosctl snapshot delete --snapshot-id SNAPSHOT_ID --wait-for-request
```

