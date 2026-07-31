---
description: "Create a Volume"
---

# VolumeCreate

## Usage

```text
ionosctl compute volume create [flags]
```

## Aliases

For `volume` command:

```text
[v vol]
```

For `create` command:

```text
[c]
```

## Description

Create a block storage Volume inside a Virtual Data Center. This does NOT attach the Volume to a Server; attach it afterwards with `ionosctl compute server volume attach`.

Storage tier (--type) determines the price/performance profile and is fixed for the life of the Volume:
  * HDD          - spinning disks, lowest cost. Best for backups, archives and cold storage. (~1,100 IOPS, up to 2,500 burst.)
  * SSD Standard - general-purpose flash. ~40 read / 30 write IOPS per GB, up to 24,000/18,000 IOPS per volume.
  * SSD Premium  - high-performance flash for databases and latency-sensitive workloads. ~75 read / 50 write IOPS per GB, up to 45,000/30,000 IOPS per volume.
  * DAS          - Direct-Attached NVMe storage that ships with Cube instances. It is bound to the Cube, its size is fixed by the Cube template, and it cannot be resized, detached or deleted independently.
Performance scales with volume size (IOPS-per-GB), so IONOS recommends booking SSD volumes of at least 100 GB. Volume sizes range from 1 GB up to 4 TB (larger on request).

Blank vs. bootable Volume:
  * Blank data disk: omit --image/--image-alias and set --licence-type so the platform knows how to bill/handle the disk. The disk is unformatted; partition and format it from the OS after attaching.
  * Bootable OS disk: pass --image (Image or Snapshot Id) OR --image-alias. When an image is set, --licence-type is derived automatically from the image and should not be overridden. For IONOS public images you must seed initial credentials with --password and/or --ssh-key-paths, otherwise you will not be able to log in. Setting --password even alongside SSH keys is recommended so the DCD remote console can authenticate with a password.

cloud-init: --user-data injects a base64-encoded cloud-init configuration on first boot (users, packages, scripts). It requires a cloud-init capable public image or image-alias.

Immutability: --type and --availability-zone are chosen at creation and cannot be changed afterwards. --size can be increased later but never decreased.

Use `--wait` (`-w`) to wait for the Volume to reach AVAILABLE state before the command returns.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string             Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -z, --availability-zone string   The storage availability zone the Volume is physically placed in. AUTO lets the platform pick; ZONE_1/ZONE_2/ZONE_3 pin it to a specific zone (useful to spread replicas across failure domains). Immutable after provisioning - to move a Volume to another zone you must snapshot it and re-create from the snapshot (default "AUTO")
      --backupunit-id string       The Id of a Backup Unit you own, used to schedule automatic backups of this Volume. Only valid on a bootable Volume, so it must be combined with a public --image or --image-alias
      --bus string                 The virtual bus the disk is exposed on once attached to a Server. VIRTIO is the paravirtualized, high-performance default and is recommended for all modern OSes. IDE is a legacy, lower-performance bus needed only in special cases (e.g. temporarily during a Windows install before VirtIO drivers are available) (default "VIRTIO")
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId]
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cpu-hot-plug               Advertise to the guest OS that CPUs can be added to the attached Server without a reboot. E.g.: --cpu-hot-plug=true. Usually inherited from the image and rarely set by hand
      --datacenter-id string       The unique Data Center Id (required)
  -D, --depth int                  Level of detail for response objects (default 1)
      --disc-virtio-hot-plug       Advertise to the guest OS that a VirtIO storage volume can be attached without a reboot. E.g.: --disc-virtio-plug=true
      --disc-virtio-hot-unplug     Advertise to the guest OS that a VirtIO storage volume can be detached without a reboot. Not supported by Windows guests. E.g.: --disc-virtio-unplug=true
  -F, --filters strings            Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
  -a, --image-alias string         A stable, human-readable alias for an IONOS public image (e.g. ubuntu:latest) to use instead of a raw --image Id. Same credential rules apply: seed --password and/or --ssh-key-paths
      --image-id string            Id of an Image or Snapshot to clone onto the Volume, making it bootable. Mutually exclusive with --image-alias. For an IONOS public image you must also seed credentials with --password and/or --ssh-key-paths. The image must match the Volume's location and disk type
      --licence-type string        The OS licence the Volume is billed and configured for. Use for blank data disks so the platform knows how to handle them. When --image or --image-alias is set, the licence type is derived automatically from the image and this flag should not be set. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER (default "LINUX")
      --limit int                  Maximum number of items to return per request (default 50)
  -n, --name string                A human-friendly label for the Volume. Not required to be unique (default "Unnamed Volume")
      --nic-hot-plug               Advertise to the guest OS that a NIC can be added without a reboot. E.g.: --nic-hot-plug=true
      --nic-hot-unplug             Advertise to the guest OS that a NIC can be removed without a reboot. E.g.: --nic-hot-unplug=true
      --no-headers                 Don't print table headers when table output is used
      --offset int                 Number of items to skip before starting to collect the results
      --order-by string            Property to order the results by
  -o, --output string              Desired output format [text|json|api-json] (default "text")
  -p, --password string            Initial root/Administrator password baked into the OS on first boot. Public images only, and immutable afterwards (rotate it from inside the guest). Allowed characters: a-z, A-Z, 0-9, 8-50 chars. Recommended even when using SSH keys so the DCD remote console can log in
      --query string               JMESPath query string to filter the output
  -q, --quiet                      Quiet output
      --ram-hot-plug               Advertise to the guest OS that memory can be added to the attached Server without a reboot. E.g.: --ram-hot-plug=true
  -s, --size string                The capacity of the Volume. Accepts a plain number (interpreted as GB) or a unit suffix, e.g. --size 10 or --size 10GB or --size 1TB. Range is 1 GB to 4 TB (larger on request); the upper bound is also capped by your contract limit. Can be increased later but never decreased (default "10")
  -k, --ssh-key-paths string       Comma-separated absolute paths to public SSH key files to authorize for the image's default user on first boot. Public images only. e.g. --ssh-key-paths "$HOME/.ssh/id_rsa.pub,/keys/ops.pub"
  -t, --timeout int                Timeout in seconds for --wait and other wait operations (default 600)
      --type string                The storage tier (fixed for the life of the Volume). HDD is cheapest (backups/cold storage); 'SSD Standard' is general-purpose flash; 'SSD Premium' is high-IOPS flash for databases; DAS is the fixed NVMe disk bundled with Cube instances (default "HDD")
      --user-data string           A base64-encoded cloud-init configuration applied on first boot (create users, install packages, run scripts). Requires a cloud-init capable public --image or --image-alias. Encode a file with e.g. base64 -w0 cloud-init.yaml
  -v, --verbose count              Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                       Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Blank 50GB SSD Premium data disk (attach to a server separately)
ionosctl compute volume create --datacenter-id DATACENTER_ID --name data-disk --size 50GB --type "SSD Premium" --licence-type LINUX

# Bootable Linux volume from a public image alias, seeded with SSH keys and a console password
ionosctl compute volume create --datacenter-id DATACENTER_ID --name boot-disk --size 20GB --type "SSD Standard" --image-alias ubuntu:latest --ssh-key-paths "$HOME/.ssh/id_rsa.pub" --password 'S3curePassw0rd!' --user-data "$(base64 -w0 cloud-init.yaml)"
```

