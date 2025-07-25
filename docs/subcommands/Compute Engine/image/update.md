---
description: "Update a specified Image"
---

# ImageUpdate

## Usage

```text
ionosctl image update [flags]
```

## Aliases

For `image` command:

```text
[img]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update information about a specified Image.

Required values to run command:

* Image Id

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --application-type string   The type of application that is hosted on this resource. Can be one of: MSSQL-2019-Web, MSSQL-2019-Standard, MSSQL-2019-Enterprise, MSSQL-2022-Web, MSSQL-2022-Standard, MSSQL-2022-Enterprise, UNKNOWN (default "UNKNOWN")
      --cloud-init string         Cloud init compatibility. Can be one of: V1, NONE (default "V1")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cpu-hot-plug              'Hot-Plug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug (default true)
      --cpu-hot-unplug            'Hot-Unplug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug
  -D, --depth int32               Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string        Description of the Image
      --disc-scsi-hot-plug        'Hot-Plug' SCSI drive (default true)
      --disc-scsi-hot-unplug      'Hot-Unplug' SCSI drive
      --disc-virtio-hot-plug      'Hot-Plug' Virt-IO drive (default true)
      --disc-virtio-hot-unplug    'Hot-Unplug' Virt-IO drive
      --expose-serial true        If set to true will expose the serial id of the disk attached to the server
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
  -i, --image-id string           The unique Image Id (required)
      --licence-type string       The OS type of this image. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER (default "UNKNOWN")
  -n, --name string               Name of the Image
      --nic-hot-plug              'Hot-Plug' NIC (default true)
      --nic-hot-unplug            'Hot-Unplug' NIC
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
      --ram-hot-plug              'Hot-Plug' RAM (default true)
      --ram-hot-unplug            'Hot-Unplug' RAM
      --require-legacy-bios       Indicates if the image requires the legacy BIOS for compatibility or specific needs. (default true)
  -t, --timeout int               Timeout option for Request for Image update [seconds] (default 60)
  -v, --verbose                   Print step-by-step process when running command
  -w, --wait-for-request          Wait for the Request for Image update to be executed
```

## Examples

```text
ionosctl image update --image-id IMAGE_ID
```

