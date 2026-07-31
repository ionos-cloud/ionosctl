---
description: "Update a Volume"
---

# VolumeUpdate

## Usage

```text
ionosctl compute volume update [flags]
```

## Aliases

For `volume` command:

```text
[v vol]
```

For `update` command:

```text
[u up]
```

## Description

Update the mutable properties of an existing Volume.

Resizing: --size may only GROW the Volume; the Cloud API cannot shrink a Volume once provisioned. If the attached Server (and the guest OS) supports disk hot-plug, the new capacity appears live without a reboot. The extra space is raw - it is NOT added to any partition or filesystem automatically, so you must extend the partition/filesystem from inside the operating system afterwards.

Immutable properties: the storage tier (--type), availability zone and the bootable image/credentials are fixed at creation and cannot be changed here. --name and --bus can be adjusted; the hot-plug capability flags advertise what the disk supports to the guest.

Use `--wait` (`-w`) to wait for the Volume to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
  -u, --api-url string           Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --bus string               The virtual bus the disk is exposed on. VIRTIO is the high-performance default; IDE is a legacy fallback (default "VIRTIO")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId]
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cpu-hot-plug             Advertise to the guest OS that CPUs can be added without a reboot. E.g.: --cpu-hot-plug=true
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int                Level of detail for response objects (default 1)
      --disc-virtio-hot-plug     Advertise to the guest OS that a VirtIO storage volume can be attached without a reboot. E.g.: --disc-virtio-plug=true
      --disc-virtio-hot-unplug   Advertise to the guest OS that a VirtIO storage volume can be detached without a reboot. Not supported by Windows guests. E.g.: --disc-virtio-unplug=true
  -F, --filters strings          Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --limit int                Maximum number of items to return per request (default 50)
  -n, --name string              A new human-friendly label for the Volume
      --nic-hot-plug             Advertise to the guest OS that a NIC can be added without a reboot. E.g.: --nic-hot-plug=true
      --nic-hot-unplug           Advertise to the guest OS that a NIC can be removed without a reboot. E.g.: --nic-hot-unplug=true
      --no-headers               Don't print table headers when table output is used
      --offset int               Number of items to skip before starting to collect the results
      --order-by string          Property to order the results by
  -o, --output string            Desired output format [text|json|api-json] (default "text")
      --query string             JMESPath query string to filter the output
  -q, --quiet                    Quiet output
      --ram-hot-plug             Advertise to the guest OS that memory can be added without a reboot. E.g.: --ram-hot-plug=true
      --size string              The new capacity of the Volume. Can only be increased, never decreased. Accepts a plain number (GB) or a unit suffix, e.g. --size 20 or --size 20GB. Upper bound is 4 TB (larger on request) and your contract limit. Remember to extend the partition/filesystem inside the guest OS afterwards
  -t, --timeout int              Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
  -i, --volume-id string         The unique Volume Id (required)
  -w, --wait                     Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Grow a volume to 20 GB (extend the filesystem inside the OS afterwards)
ionosctl compute volume update --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --size 20GB

# Rename a volume and enable RAM hot-plug advertisement
ionosctl compute volume update --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --name prod-data --ram-hot-plug=true
```

