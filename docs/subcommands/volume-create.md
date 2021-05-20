---
description: Create a Volume
---

# VolumeCreate

## Usage

```text
ionosctl volume create [flags]
```

## Aliases

For `volume` command:
```text
[v]
```

## Description

Use this command to create a Volume on your account. Creates a volume within the Data Center. This will NOT attach the Volume to a Server. Please see the Servers commands for details on how to attach storage Volumes. You can specify the name, size, type, licence type, availability zone, image and other properties for the object.

Note: You will need to provide a valid value for either the Image, Image Alias, or the Licence Type options. The Licence Type is required, but if Image or Image Alias is supplied, then Licence Type is already set and cannot be changed. Similarly either the Image Password or SSH Keys attributes need to be defined when creating a Volume that uses an Image or Image Alias of an IONOS public HDD Image. You may wish to set a valid value for Image Password even when using SSH Keys so that it is possible to authenticate with a password when using the remote console feature of the DCD.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Licence Type/Image Id or Image Alias

## Options

```text
  -u, --api-url string             Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -z, --availability-zone string   Availability zone of the Volume. Storage zone can only be selected prior provisioning (default "AUTO")
      --backupunit-id string       The unique Id of the Backup Unit that User has access to. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property
      --bus string                 Bus for the Volume (default "VIRTIO")
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId SshKeys ImageAlias DeviceNumber UserData] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cpu-hot-plug               It is capable of CPU hot plug (no reboot required)
      --datacenter-id string       The unique Data Center Id (required)
      --disc-virtio-hot-plug       It is capable of Virt-IO drive hot plug (no reboot required)
      --disc-virtio-hot-unplug     It is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines
  -f, --force                      Force command to execute without user input
  -h, --help                       help for create
      --image-alias string         The Image Alias to set instead of Image Id
      --image-id string            The Image Id or snapshot Id to be used as template for the new Volume
      --licence-type string        Licence Type of the Volume
  -n, --name string                Name of the Volume
      --nic-hot-plug               It is capable of nic hot plug (no reboot required)
      --nic-hot-unplug             It is capable of nic hot unplug (no reboot required)
  -o, --output string              Desired output format [text|json] (default "text")
  -p, --password string            Initial password to be set for installed OS. Works with public Images only. Not modifiable. Password rules allows all characters from a-z, A-Z, 0-9
  -q, --quiet                      Quiet output
      --ram-hot-plug               It is capable of memory hot plug (no reboot required)
      --size float32               Size in GB of the Volume (default 10)
      --ssh-keys strings           SSH Keys of the Volume
  -t, --timeout int                Timeout option for Request for Volume creation [seconds] (default 60)
      --type string                Type of the Volume (default "HDD")
      --user-data string           The cloud-init configuration for the Volume as base64 encoded string. It is mandatory to provide either 'public image' or 'imageAlias' that has cloud-init compatibility in conjunction with this property
  -w, --wait-for-request           Wait for the Request for Volume creation to be executed
```

## Examples

```text
ionosctl volume create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --name demoVolume
VolumeId                               Name         Size   Type   LicenceType   State   Image
ce510144-9bc6-4115-bd3d-b9cd232dd422   demoVolume   10GB   HDD    LINUX         BUSY    
RequestId: a2da3bb7-3851-4e80-a5e9-6e98a66cebab
Status: Command volume create has been successfully executed
```

