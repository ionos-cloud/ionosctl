---
description: "Upload an image to FTP server"
---

# ImageUpload

## Usage

```text
ionosctl image upload [flags]
```

## Aliases

For `image` command:

```text
[img]
```

For `upload` command:

```text
[ftp-upload ftp upl]
```

## Description

Use this command to upload an HDD or ISO image.

Required values to run command:

* Location


## Options

```text
  -u, --api-url string           Override default host url (default "https://api.ionos.com")
      --cloud-init string        Cloud init compatibility. Can be one of: V1, NONE (default "V1")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cpu-hot-plug             'Hot-Plug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug (default true)
      --cpu-hot-unplug           'Hot-Unplug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug
      --crt-path string          (Unneeded for IONOS FTP Servers) Path to file containing server certificate. If your FTP server is self-signed, you need to add the server certificate to the list of certificate authorities trusted by the client.
  -d, --description string       Description of the Image
      --disc-scsi-hot-plug       'Hot-Plug' SCSI drive (default true)
      --disc-scsi-hot-unplug     'Hot-Unplug' SCSI drive
      --disc-virtio-hot-plug     'Hot-Plug' Virt-IO drive (default true)
      --disc-virtio-hot-unplug   'Hot-Unplug' Virt-IO drive
      --ftp-url string           URL of FTP server, with %s flag if location is embedded into url (default "ftp-%s.ionos.com")
  -h, --help                     Print usage
  -i, --image strings            Slice of paths to images, can be absolute path or relative to current working directory (required)
  -a, --image-alias strings      Rename the uploaded images. These names should not contain any extension. By default, this is the base of the image path
      --licence-type string      The OS type of this image. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, UNKNOWN, OTHER (default "UNKNOWN")
  -l, --location strings         Location to upload to. Must be an array containing only fra, fkb, txl, lhr, las, ewr, vit (required)
  -n, --name string              Name of the Image
      --nic-hot-plug             'Hot-Plug' NIC (default true)
      --nic-hot-unplug           'Hot-Unplug' NIC
      --no-headers               Don't print table headers when table output is used
  -o, --output string            Desired output format [text|json|api-json] (default "text")
  -q, --quiet                    Quiet output
      --ram-hot-plug             'Hot-Plug' RAM (default true)
      --ram-hot-unplug           'Hot-Unplug' RAM
      --skip-update              After the image is uploaded to the FTP server, send a PATCH to the API with the contents of the image properties flags and emulate a "create" command.
      --skip-verify              Skip verification of server certificate, useful if using a custom ftp-url. WARNING: You can be the target of a man-in-the-middle attack!
  -t, --timeout int              (seconds) Context Deadline. FTP connection will time out after this many seconds (default 300)
  -v, --verbose                  Print step-by-step process when running command
```

## Examples

```text
ionosctl img u -i kolibri.iso -l fkb,fra,vit
```

