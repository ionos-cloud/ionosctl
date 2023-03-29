---
description: Get a Snapshot
---

# SnapshotGet

## Usage

```text
ionosctl snapshot get [flags]
```

## Aliases

For `snapshot` command:

```text
[ss snap]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Snapshot.

Required values to run command:

* Snapshot Id

## Options

```text
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
      --no-headers           When using text output, don't print headers
  -i, --snapshot-id string   The unique Snapshot Id (required)
```

## Examples

```text
ionosctl snapshot get --snapshot-id SNAPSHOT_ID
```

