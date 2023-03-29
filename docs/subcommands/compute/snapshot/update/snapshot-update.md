---
description: Update a Snapshot
---

# SnapshotUpdate

## Usage

```text
ionosctl snapshot update [flags]
```

## Aliases

For `snapshot` command:

```text
[ss snap]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified Snapshot.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Snapshot Id

## Options

```text
      --cpu-hot-plug             This volume is capable of CPU hot plug (no reboot required). E.g.: --cpu-hot-plug=true, --cpu-hot-plug=false
      --cpu-hot-unplug           This volume is capable of CPU hot unplug (no reboot required). E.g.: --cpu-hot-unplug=true, --cpu-hot-unplug=false
  -D, --depth int32              Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string       Description of the Snapshot
      --disc-scsi-hot-plug       This volume is capable of SCSI drive hot plug (no reboot required). E.g.: --disc-scsi-plug=true, --disc-scsi-plug=false
      --disc-scsi-hot-unplug     This volume is capable of SCSI drive hot unplug (no reboot required). E.g.: --disc-scsi-unplug=true, --disc-scsi-unplug=false
      --disc-virtio-hot-plug     This volume is capable of VirtIO drive hot plug (no reboot required). E.g.: --disc-virtio-plug=true, --disc-virtio-plug=false
      --disc-virtio-hot-unplug   This volume is capable of VirtIO drive hot unplug (no reboot required). E.g.: --disc-virtio-unplug=true, --disc-virtio-unplug=false
      --licence-type string      Licence Type of the Snapshot
  -n, --name string              Name of the Snapshot
      --nic-hot-plug             This volume is capable of NIC hot plug (no reboot required). E.g.: --nic-hot-plug=true, --nic-hot-plug=false
      --nic-hot-unplug           This volume is capable of NIC hot unplug (no reboot required). E.g.: --nic-hot-unplug=true, --nic-hot-unplug=false
      --ram-hot-plug             This volume is capable of memory hot plug (no reboot required). E.g.: --ram-hot-plug=true, --ram-hot-plug=false
      --ram-hot-unplug           This volume is capable of memory hot unplug (no reboot required). E.g.: --ram-hot-unplug=true, --ram-hot-unplug=false
      --sec-auth-protection      Enable secure authentication protection. E.g.: --sec-auth-protection=true, --sec-auth-protection=false
  -i, --snapshot-id string       The unique Snapshot Id (required)
  -t, --timeout int              Timeout option for Request for Snapshot creation [seconds] (default 60)
  -w, --wait-for-request         Wait for the Request for Snapshot creation to be executed
```

## Examples

```text
ionosctl snapshot update --snapshot-id SNAPSHOT_ID --name NAME
```

