---
description: "Update a Snapshot's metadata and inherited capabilities"
---

# SnapshotUpdate

## Usage

```text
ionosctl compute snapshot update [flags]
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

Use this command to update the metadata of an existing Snapshot. This edits properties recorded ON the Snapshot - it does not re-capture or change the stored data image.

Two kinds of properties can be changed:
  1. Descriptive: --name, --description, --licence-type, --sec-auth-protection.
  2. Hot-plug capability hints (cpu/ram/nic/disc SCSI/disc VirtIO, plug and unplug). These flags advertise whether a component can be added/removed at runtime without a reboot. A Volume created or restored from this Snapshot inherits these values, so setting them here fixes up the defaults future Volumes will get. Only pass the flags you want to change; unspecified capabilities are left untouched. VirtIO disk hot-unplug is unsupported on Windows, and SCSI disk hot-unplug is limited to non-Windows guests.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Snapshot Id

## Options

```text
  -u, --api-url string           Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [SnapshotId Name LicenceType Size State]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cpu-hot-plug             Advertise that Volumes from this Snapshot support adding vCPUs at runtime without a reboot. E.g.: --cpu-hot-plug=true, --cpu-hot-plug=false
      --cpu-hot-unplug           Advertise that Volumes from this Snapshot support removing vCPUs at runtime without a reboot. E.g.: --cpu-hot-unplug=true, --cpu-hot-unplug=false
  -D, --depth int                Level of detail for response objects (default 1)
  -d, --description string       Free-form notes about the Snapshot, e.g. why or when it was taken
      --disc-scsi-hot-plug       Advertise that Volumes from this Snapshot support attaching a SCSI disk at runtime without a reboot. E.g.: --disc-scsi-plug=true, --disc-scsi-plug=false
      --disc-scsi-hot-unplug     Advertise that Volumes from this Snapshot support detaching a SCSI disk at runtime without a reboot. Limited to non-Windows guests. E.g.: --disc-scsi-unplug=true, --disc-scsi-unplug=false
      --disc-virtio-hot-plug     Advertise that Volumes from this Snapshot support attaching a VirtIO disk at runtime without a reboot. E.g.: --disc-virtio-plug=true, --disc-virtio-plug=false
      --disc-virtio-hot-unplug   Advertise that Volumes from this Snapshot support detaching a VirtIO disk at runtime without a reboot. Not supported on Windows guests. E.g.: --disc-virtio-unplug=true, --disc-virtio-unplug=false
  -F, --filters strings          Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --licence-type string      The operating-system licence recorded on the Snapshot. Inherited by Volumes created/restored from it and affects OS licensing (WINDOWS variants are billed). Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER
      --limit int                Maximum number of items to return per request (default 50)
  -n, --name string              A human-friendly label for the Snapshot; shown in listings. Does not have to be unique
      --nic-hot-plug             Advertise that Volumes from this Snapshot support attaching a NIC at runtime without a reboot. E.g.: --nic-hot-plug=true, --nic-hot-plug=false
      --nic-hot-unplug           Advertise that Volumes from this Snapshot support detaching a NIC at runtime without a reboot. E.g.: --nic-hot-unplug=true, --nic-hot-unplug=false
      --no-headers               Don't print table headers when table output is used
      --offset int               Number of items to skip before starting to collect the results
      --order-by string          Property to order the results by
  -o, --output string            Desired output format [text|json|api-json] (default "text")
      --query string             JMESPath query string to filter the output
  -q, --quiet                    Quiet output
      --ram-hot-plug             Advertise that Volumes from this Snapshot support adding memory at runtime without a reboot. E.g.: --ram-hot-plug=true, --ram-hot-plug=false
      --ram-hot-unplug           Advertise that Volumes from this Snapshot support removing memory at runtime without a reboot. E.g.: --ram-hot-unplug=true, --ram-hot-unplug=false
      --sec-auth-protection      Protect the Snapshot with secure authentication: when true, deleting or restoring it requires the Contract Owner or a re-authenticated user. E.g.: --sec-auth-protection=true, --sec-auth-protection=false
  -i, --snapshot-id string       The unique Snapshot Id (required)
  -t, --timeout int              Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                     Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a snapshot and fix its OS licence
ionosctl compute snapshot update --snapshot-id SNAPSHOT_ID --name "prod-db golden v2" --licence-type LINUX

# Advanced: mark the image as CPU/RAM/NIC hot-plug capable so volumes restored from it inherit those capabilities
ionosctl compute snapshot update --snapshot-id SNAPSHOT_ID --cpu-hot-plug=true --ram-hot-plug=true --nic-hot-plug=true --disc-virtio-plug=true --wait
```

