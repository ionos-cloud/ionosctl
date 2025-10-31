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

This command uploads one or more disk images to an FTP server using FTP over TLS (FTPS), then optionally updates the uploaded images via the Images API to set properties you passed as flags.
This command requires that you are logged in using IONOS_USERNAME and IONOS_PASSWORD environment variables.

High level steps:
  1. Upload file(s) concurrently to the target FTP server(s).
  2. If you do not use --skip-update, poll the Images API for the uploaded image(s) to appear.
  3. When the API shows the uploaded image(s), perform PATCH /images/<UUID> to apply the requested image properties.
  4. Print the resulting image objects to stdout in the chosen table or JSON format.

AUTH AND SAFETY
  - The FTP server relies on API credentials via environment variables IONOS_USERNAME and IONOS_PASSWORD. You can debug your current setup with "ionosctl whoami --provenance".
  - Use --skip-update to skip the API PATCH step if you only want to perform an FTP upload and not modify images through the API.
  - Use --skip-verify to skip verifying the FTP server certificate. Only use that for trusted servers. Skipping certificate verification can expose you to man-in-the-middle attacks.
  - If using a custom FTP server it is advised to use a self-signed certificate instead of --skip-verify. Provide its PEM file via --crt-path. The file should contain the server certificate in base64 PEM format.

FTP URLs
  - Default IONOS FTP servers are of the form ftp-<region>.ionos.com (for example ftp-fra.ionos.com).
  - If uploading to default IONOS FTP servers, --ftp-url is optional. The command will construct the URL automatically from the locations you provide via --location (i.e. 'de/fra' or 'fra').
  - The command chooses the remote path automatically:
      * Files ending in .iso or .img are uploaded to the iso-images/ directory.
      * All other supported image extensions are uploaded to the hdd-images/ directory.
  - If you supply a custom --ftp-url that contains a placeholder, for example ftp://myftp.example/locations/%s, you must also supply one or more --location values. The command will replace %s with the location-specific fragment for each location. Example: --ftp-url ftp://myftp.example/locations/%s --location fra,fkb
  - If you supply a custom --ftp-url without a placeholder, you may provide multiple --ftp-url values to try multiple servers.

POLLING AND TIMEOUTS
  - After upload, unless you use --skip-update, the command repeatedly queries GET /images with filters for the uploaded file names and locations.
  - Polling runs until either all expected images appear, or the command context deadline expires.
  - The context deadline is controlled with --timeout (seconds). The FTP connection and the subsequent API operations share the same context. If a timeout occurs after FTP finished but before the PATCH completed, the PATCH will be cancelled.

NOTES
  - Uploading multiple images with the same name to the same location is forbidden.
  - The command does not delete or overwrite existing images on the FTP server. If an image with the same name already exists on the server, the upload will fail.
  - The command does not check if the uploaded image is valid or bootable. It only checks the file extension.
  - You can use 'ionosctl image list --filter public=false' to see your uploaded images.
  - You must contact support to delete images you uploaded via FTP. Deleting them via API will only set their size to 0B.

EXAMPLES
  - Simple upload to IONOS servers:
    ionosctl img upload -i image.iso -l de/fra,de/fkb,es/vit --skip-update
    Uploads image.iso to ftp://ftp-fkb.ionos.com/iso-images, ftp://ftp-fra.ionos.com/iso-images and ftp://ftp-vit.ionos.com/iso-images, then exits without calling the Images API.

  - Upload and let the CLI set properties via API:
    ionosctl img upload -i image.iso -l de/fra
    Uploads to ftp://ftp-fra.ionos.com/iso-images, polls GET /images until the image appears, then PATCHes that image with the properties you supplied via flags.

  - Use a custom FTP server:
    ionosctl img upload -i image.iso --ftp-url "ftp://myftp.example" --crt-path certificates/my-server-crt.pem --skip-update

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
      --limit int                 Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location strings          Location to upload to. Can be one of de/fra, de/fra/2, es/vit, gb/lhr, gb/bhx, fr/par, us/las, us/ewr, us/mci, de/txl, de/fkb if not using --ftp-url (required)
  -n, --name string               Name of the Image
      --nic-hot-plug              'Hot-Plug' NIC (default true)
      --nic-hot-unplug            'Hot-Unplug' NIC
      --no-headers                Don't print table headers when table output is used
      --offset int                Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
      --ram-hot-plug              'Hot-Plug' RAM (default true)
      --ram-hot-unplug            'Hot-Unplug' RAM
      --rename strings            Rename the uploaded images before trying to upload. These names should not contain any extension. By default, this is the base of the image path
      --require-legacy-bios       Indicates if the image requires the legacy BIOS for compatibility or specific needs. (default true)
      --skip-update               Skip setting image properties after it has been uploaded. Normal behavior is to send a PATCH to the API, after the image has been uploaded, with the contents of the image properties flags and emulate a "create" command.
      --skip-verify               Skip verification of server certificate, useful if using a custom ftp-url. WARNING: You can be the target of a man-in-the-middle attack!
  -t, --timeout int               (seconds) Context Deadline. FTP connection will time out after this many seconds (default 1200)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
```

