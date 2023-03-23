---
description: Restore a Snapshot onto a Volume
---

# SnapshotRestore

## Usage

```text
ionosctl snapshot restore [flags]
```

## Aliases

For `snapshot` command:

```text
[ss snap]
```

For `restore` command:

```text
[r]
```

## Description

Use this command to restore a Snapshot onto a Volume. A Snapshot is created as just another image that can be used to create new Volumes or to restore an existing Volume.

Required values to run command:

* Datacenter Id
* Volume Id
* Snapshot Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -i, --snapshot-id string     The unique Snapshot Id (required)
  -t, --timeout int            Timeout option for Request for Snapshot restore [seconds] (default 60)
      --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Snapshot restore to be executed
```

## Examples

```text
ionosctl snapshot restore --snapshot-id SNAPSHOT_ID --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --wait-for-request
```

