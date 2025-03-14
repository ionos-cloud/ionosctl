---
description: "Update a Volume"
---

# VolumeUpdate

## Usage

```text
ionosctl volume update [flags]
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

Use this command to update a Volume. You may increase the size of an existing storage Volume. You cannot reduce the size of an existing storage Volume. The Volume size will be increased without reboot if the appropriate "hot plug" settings have been set to true. The additional capacity is not added to any partition therefore you will need to adjust the partition inside the operating system afterwards.

Once you have increased the Volume size you cannot decrease the Volume size using the Cloud API. Certain attributes can only be set when a Volume is created and are considered immutable once the Volume has been provisioned.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
  -u, --api-url string           Override default host url (default "https://api.ionos.com")
      --bus string               Bus of the Volume (default "VIRTIO")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cpu-hot-plug             It is capable of CPU hot plug (no reboot required). E.g.: --cpu-hot-plug=true, --cpu-hot-plug=false
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int32              Controls the detail depth of the response objects. Max depth is 10.
      --disc-virtio-hot-plug     It is capable of Virt-IO drive hot plug (no reboot required). E.g.: --disc-virtio-plug=true, --disc-virtio-plug=false
      --disc-virtio-hot-unplug   It is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines. E.g.: --disc-virtio-unplug=true, --disc-virtio-unplug=false
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
  -n, --name string              Name of the Volume
      --nic-hot-plug             It is capable of nic hot plug (no reboot required). E.g.: --nic-hot-plug=true, --nic-hot-plug=false
      --nic-hot-unplug           It is capable of nic hot unplug (no reboot required). E.g.: --nic-hot-unplug=true, --nic-hot-unplug=false
      --no-headers               Don't print table headers when table output is used
  -o, --output string            Desired output format [text|json|api-json] (default "text")
  -q, --quiet                    Quiet output
      --ram-hot-plug             It is capable of memory hot plug (no reboot required). E.g.: --ram-hot-plug=true, --ram-hot-plug=false
      --size string              The size of the Volume in GB. e.g. 10 or 10GB. The maximum volume size is determined by your contract limit
  -t, --timeout int              Timeout option for Request for Volume update [seconds] (default 60)
  -v, --verbose                  Print step-by-step process when running command
  -i, --volume-id string         The unique Volume Id (required)
  -w, --wait                     Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl volume update --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --size 20
```

