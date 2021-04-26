---
description: Create a Volume
---

# Create

## Usage

```text
ionosctl volume create [flags]
```

## Description

Use this command to create a Volume on your account. You can specify the name, size, type, licence type and availability zone for the object.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string               Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings                 Columns to be printed in the standard output (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string         The unique Data Center Id [Required flag]
      --force                        Force command to execute without user input
  -h, --help                         help for create
  -o, --output string                Desired output format [text|json] (default "text")
  -q, --quiet                        Quiet output
      --timeout int                  Timeout option for Volume to be created [seconds] (default 60)
      --volume-bus string            Bus for the Volume (default "VIRTIO")
      --volume-licence-type string   Licence Type of the Volume (default "LINUX")
      --volume-name string           Name of the Volume
      --volume-size float32          Size in GB of the Volume (default 10)
      --volume-ssh-keys string       Ssh Key of the Volume
      --volume-type string           Type of the Volume (default "HDD")
      --volume-zone string           Availability zone of the Volume. Storage zone can only be selected prior provisioning (default "AUTO")
      --wait                         Wait for Volume to be created
```

## Examples

```text
ionosctl volume create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --volume-name demoVolume
VolumeId                               Name         Size   Type   LicenceType   State   Image
ce510144-9bc6-4115-bd3d-b9cd232dd422   demoVolume   10GB   HDD    LINUX         BUSY    
RequestId: a2da3bb7-3851-4e80-a5e9-6e98a66cebab
Status: Command volume create has been successfully executed
```

