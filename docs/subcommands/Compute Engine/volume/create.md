---
description: "Create a Volume"
---

# VolumeCreate

## Usage

```text
ionosctl volume create [flags]
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

Use this command to create a Volume on your account, within a Data Center. This will NOT attach the Volume to a Server. Please see the Servers commands for details on how to attach storage Volumes. You can specify the name, size, type, licence type, availability zone, image and other properties for the object.

NNote: The Licence Type has a default value, but if Image ID or Image Alias is supplied, then Licence Type will be automatically set. The Image Password or SSH Keys attributes can be defined when creating a Volume that uses an Image ID or Image Alias of an IONOS public Image. You may wish to set a valid value for Image Password even when using SSH Keys so that it is possible to authenticate with a password when using the remote console feature of the DCD.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string             Override default host url (default "https://api.ionos.com")
  -z, --availability-zone string   Availability zone of the Volume. Storage zone can only be selected prior provisioning (default "AUTO")
      --backupunit-id string       The unique Id of the Backup Unit that User has access to. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property
      --bus string                 The bus type of the Volume (default "VIRTIO")
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [VolumeId Name Size Type LicenceType State Image Bus AvailabilityZone BackupunitId DeviceNumber UserData BootServerId DatacenterId] (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cpu-hot-plug               It is capable of CPU hot plug (no reboot required). E.g.: --cpu-hot-plug=true, --cpu-hot-plug=false
      --datacenter-id string       The unique Data Center Id (required)
  -D, --depth int32                Controls the detail depth of the response objects. Max depth is 10.
      --disc-virtio-hot-plug       It is capable of Virt-IO drive hot plug (no reboot required). E.g.: --disc-virtio-plug=true, --disc-virtio-plug=false
      --disc-virtio-hot-unplug     It is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines. E.g.: --disc-virtio-unplug=true, --disc-virtio-unplug=false
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
  -a, --image-alias string         The Image Alias to set instead of Image Id. A password or SSH Key need to be set
      --image-id string            The Image Id or Snapshot Id to be used as template for the new Volume. A password or SSH Key need to be set
      --licence-type string        Licence Type of the Volume. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, UNKNOWN, OTHER (default "LINUX")
  -n, --name string                Name of the Volume (default "Unnamed Volume")
      --nic-hot-plug               It is capable of nic hot plug (no reboot required). E.g.: --nic-hot-plug=true, --nic-hot-plug=false
      --nic-hot-unplug             It is capable of nic hot unplug (no reboot required). E.g.: --nic-hot-unplug=true, --nic-hot-unplug=false
  -o, --output string              Desired output format [text|json|api-json] (default "text")
  -p, --password string            Initial password to be set for installed OS. Works with public Images only. Not modifiable. Password rules allows all characters from a-z, A-Z, 0-9
  -q, --quiet                      Quiet output
      --ram-hot-plug               It is capable of memory hot plug (no reboot required). E.g.: --ram-hot-plug=true, --ram-hot-plug=false
  -s, --size string                The size of the Volume in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit (default "10")
  -k, --ssh-key-paths string       Absolute paths of the SSH Keys for the Volume
  -t, --timeout int                Timeout option for Request for Volume creation [seconds] (default 60)
      --type string                Type of the Volume (default "HDD")
      --user-data string           The cloud-init configuration for the Volume as base64 encoded string. It is mandatory to provide either 'public image' or 'imageAlias' that has cloud-init compatibility in conjunction with this property
  -v, --verbose                    Print step-by-step process when running command
  -w, --wait-for-request           Wait for the Request for Volume creation to be executed
```

## Examples

```text
ionosctl volume create --datacenter-id DATACENTER_ID --name NAME

ionosctl volume create --datacenter-id DATACENTER_ID --name NAME --image-alias IMAGE_ALIAS --ssh-keys-path "SSH_KEY_PATH1,SSH_KEY_PATH2"
```

