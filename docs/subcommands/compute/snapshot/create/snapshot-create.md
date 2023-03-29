---
description: Create a Snapshot of a Volume within the Virtual Data Center
---

# SnapshotCreate

## Usage

```text
ionosctl snapshot create [flags]
```

## Aliases

For `snapshot` command:

```text
[ss snap]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a Snapshot. Creation of Snapshots is performed from the perspective of the storage Volume. The name, description and licence type of the Snapshot can be set.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string     Description of the Snapshot
      --licence-type string    Licence Type of the Snapshot (default "LINUX")
  -n, --name string            Name of the Snapshot (default "Unnamed Snapshot")
      --sec-auth-protection    Enable secure authentication protection. E.g.: --sec-auth-protection=true, --sec-auth-protection=false
  -t, --timeout int            Timeout option for Request for Snapshot creation [seconds] (default 60)
      --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Snapshot creation to be executed
```

## Examples

```text
ionosctl snapshot create --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --name NAME
```

