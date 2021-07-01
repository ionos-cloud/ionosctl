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
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [SnapshotId Name LicenceType Size State] (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cpu-hot-plug             This volume is capable of CPU hot plug (no reboot required)
      --cpu-hot-unplug           This volume is capable of CPU hot unplug (no reboot required)
  -d, --description string       Description of the Snapshot
      --disc-scsi-hot-plug       This volume is capable of SCSI drive hot plug (no reboot required)
      --disc-scsi-hot-unplug     This volume is capable of SCSI drive hot unplug (no reboot required)
      --disc-virtio-hot-plug     This volume is capable of VirtIO drive hot plug (no reboot required)
      --disc-virtio-hot-unplug   This volume is capable of VirtIO drive hot unplug (no reboot required)
  -f, --force                    Force command to execute without user input
  -h, --help                     help for update
      --licence-type string      Licence Type of the Snapshot
  -n, --name string              Name of the Snapshot
      --nic-hot-plug             This volume is capable of NIC hot plug (no reboot required)
      --nic-hot-unplug           This volume is capable of NIC hot unplug (no reboot required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --ram-hot-plug             This volume is capable of memory hot plug (no reboot required)
      --ram-hot-unplug           This volume is capable of memory hot unplug (no reboot required)
      --sec-auth-protection      Enable secure authentication protection
  -i, --snapshot-id string       The unique Snapshot Id (required)
  -t, --timeout int              Timeout option for Request for Snapshot creation [seconds] (default 60)
  -w, --wait-for-request         Wait for the Request for Snapshot creation to be executed
```

## Examples

```text
ionosctl snapshot update --snapshot-id SNAPSHOT_ID --name NAME
```

