---
description: Update a Volume
---

# Update

## Usage

```text
ionosctl volume update [flags]
```

## Description

Use this command to update a Volume. You may increase the size of an existing storage Volume. You cannot reduce the size of an existing storage Volume. The Volume size will be increased without reboot if the appropriate "hot plug" settings have been set to true. The additional capacity is not added to any partition therefore you will need to adjust the partition inside the operating system afterwards.

Once you have increased the Volume size you cannot decrease the Volume size using the Cloud API. Certain attributes can only be set when a Volume is created and are considered immutable once the Volume has been provisioned.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
      --force                    Force command to execute without user input
  -h, --help                     help for update
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --timeout int              Timeout option for Volume to be updated [seconds] (default 60)
      --volume-bus string        Bus of the Volume (default "VIRTIO")
      --volume-id string         The unique Volume Id (required)
      --volume-name string       Name of the Volume
      --volume-size float32      Size in GB of the Volume (default 10)
      --volume-ssh-keys string   Ssh Key of the Volume
      --wait                     Wait for Volume to be updated
```

## Examples

```text
ionosctl volume update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --volume-id ce510144-9bc6-4115-bd3d-b9cd232dd422 --volume-size 20
VolumeId                               Name         Size   Type   LicenceType   State   Image
ce510144-9bc6-4115-bd3d-b9cd232dd422   demoVolume   20GB   HDD    LINUX         BUSY    
RequestId: ad4080a9-a51f-4d81-ae40-660cbfe009f4
Status: Command volume update has been successfully executed
```

