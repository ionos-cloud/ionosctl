---
description: "Upload an image to FTP server using FTP over TLS (FTPS)"
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

WARNING:
This command can only be used if 2-Factor Authentication is disabled on your account and you're logged in using IONOS_USERNAME and IONOS_PASSWORD environment variables (see "Authenticating with Ionos Cloud" at https://docs.ionos.com/cli-ionosctl).

OVERVIEW:
  Use this command to securely upload one or more HDD or ISO images to the specified FTP server using FTP over TLS (FTPS). This command supports a variety of options to provide flexibility during the upload process:
  - The command supports renaming the uploaded images with the '--rename' flag. If uploading multiple images, you must provide a new name for each image.
  - Specify the context deadline for the FTP connection using the '--timeout' flag. The operation as a whole will terminate after the specified number of seconds, i.e. if the FTP upload had finished but your PATCH operation did not, only the PATCH operation will be intrerrupted.
POST-UPLOAD OPERATIONS:
  By default, this command will query 'GET /images' endpoint for your uploaded images, then try to use 'PATCH /images/<UUID>' to update the uploaded images with the given property fields.
  - It is necessary to use valid API credentials for this.
  - To skip this API behaviour, you can use '--skip-update'.
CUSTOM URLs:
  This command supports usage of other FTP servers too, not just the IONOS ones.
  - The '--location' flag is only required if your '--ftp-url' contains a placeholder variable (i.e. %s).
  In this case, for every location in that slice, an attempt of FTP upload would be made at the URL computed by embedding it into the placeholder variable
  - Use the '--skip-verify' flag to skip the verification of the server certificate. This can be useful when using a custom ftp-url,
  but be warned that this could expose you to a man-in-the-middle attack.
  - If you're using a self-signed FTP server, you can provide the path to the server certificate file in base64 PEM format using the '--crt-path' flag.


## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
      --application-type string   The type of application that is hosted on this resource. Can be one of: MSSQL-2019-Web, MSSQL-2019-Standard, MSSQL-2019-Enterprise, MSSQL-2022-Web, MSSQL-2022-Standard, MSSQL-2022-Enterprise, UNKNOWN (default "UNKNOWN")
      --cloud-init string         Cloud init compatibility. Can be one of: V1, NONE (default "V1")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cpu-hot-plug              'Hot-Plug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug (default true)
      --cpu-hot-unplug            'Hot-Unplug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug
      --crt-path string           (Not needed for IONOS FTP Servers) Path to file containing server certificate. If your FTP server is self-signed, you need to add the server certificate to the list of certificate authorities trusted by the client.
  -d, --description string        Description of the Image
      --disc-scsi-hot-plug        'Hot-Plug' SCSI drive (default true)
      --disc-scsi-hot-unplug      'Hot-Unplug' SCSI drive
      --disc-virtio-hot-plug      'Hot-Plug' Virt-IO drive (default true)
      --disc-virtio-hot-unplug    'Hot-Unplug' Virt-IO drive
      --expose-serial true        If set to true will expose the serial id of the disk attached to the server
  -f, --force                     Force command to execute without user input
      --ftp-url string            URL of FTP server, with %s flag if location is embedded into url (default "ftp-%s.ionos.com")
  -h, --help                      Print usage
  -i, --image strings             Slice of paths to images, can be absolute path or relative to current working directory (required)
      --licence-type string       The OS type of this image. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER (default "UNKNOWN")
  -l, --location strings          Location to upload to. Must be an array containing only fra, fkb, txl, lhr, las, ewr, vit if not using --ftp-url (required)
  -n, --name string               Name of the Image
      --nic-hot-plug              'Hot-Plug' NIC (default true)
      --nic-hot-unplug            'Hot-Unplug' NIC
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
      --ram-hot-plug              'Hot-Plug' RAM (default true)
      --ram-hot-unplug            'Hot-Unplug' RAM
      --rename strings            Rename the uploaded images before trying to upload. These names should not contain any extension. By default, this is the base of the image path
      --require-legacy-bios       Indicates if the image requires the legacy BIOS for compatibility or specific needs. (default true)
      --skip-update               Skip setting image properties after it has been uploaded. Normal behavior is to send a PATCH to the API, after the image has been uploaded, with the contents of the image properties flags and emulate a "create" command.
      --skip-verify               Skip verification of server certificate, useful if using a custom ftp-url. WARNING: You can be the target of a man-in-the-middle attack!
  -t, --timeout int               (seconds) Context Deadline. FTP connection will time out after this many seconds (default 300)
  -v, --verbose                   Print step-by-step process when running command
```

## Examples

```text
- 'ionosctl img upload -i kolibri.iso -l fkb,fra,vit --skip-update': Simply upload the image 'kolibri.iso' from the current directory to IONOS FTP servers 'ftp://ftp-fkb.ionos.com/iso-images', 'ftp://ftp-fra.ionos.com/iso-images', 'ftp://ftp-vit.ionos.com/iso-images'.
- 'ionosctl img upload -i kolibri.iso -l fra': Upload the image 'kolibri.iso' from the current directory to IONOS FTP server 'ftp://ftp-fra.ionos.com/iso-images'. Once the upload has finished, start querying 'GET /images' with a filter for 'kolibri', to get the UUID of the image as seen by the Images API. When UUID is found, perform a 'PATCH /images/<UUID>' to set the default flag values.
- 'ionosctl img upload -i kolibri.iso --skip-update --skip-verify --ftp-url ftp://12.34.56.78': Use your own custom server. Use skip verify to skip checking server's identity
- 'ionosctl img upload -i kolibri.iso -l fra --ftp-url ftp://myComplexFTPServer/locations/%s --crt-path certificates/my-servers-cert.crt --location Paris,Berlin,LA,ZZZ --skip-update': Upload the image to multiple FTP servers, with location embedding into URL.
```

