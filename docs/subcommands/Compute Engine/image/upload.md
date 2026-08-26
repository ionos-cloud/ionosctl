---
description: "Upload a local disk image to IONOS as a private image via FTP over TLS (FTPS)"
---

# ImageUpload

## Usage

```text
ionosctl compute image upload [flags]
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
  - The FTP server relies on basic API credentials via environment variables IONOS_USERNAME and IONOS_PASSWORD. A bearer token (IONOS_TOKEN) cannot be used for the FTP upload, so if you authenticate with a token you must additionally set IONOS_USERNAME and IONOS_PASSWORD (they may be set alongside IONOS_TOKEN). You can debug your current setup with "ionosctl whoami --provenance".
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
  - You can use 'ionosctl compute image list --filter public=false' to see your uploaded images.
  - You must contact support to delete images you uploaded via FTP. Deleting them via API will only set their size to 0B.

EXAMPLES
  - Simple upload to IONOS servers:
    ionosctl img upload -i image.iso -l de/fra,de/fkb,es/vit --skip-update
    Uploads image.iso to ftp://ftp-fkb.ionos.com/iso-images, ftp://ftp-fra.ionos.com/iso-images and ftp://ftp-vit.ionos.com/iso-images, then exits without calling the Images API.

  - Upload and let the CLI set properties via API (BASIC):
    ionosctl img upload -i image.iso -l de/fra
    Uploads to ftp://ftp-fra.ionos.com/iso-images, polls GET /images until the image appears, then PATCHes that image with the properties you supplied via flags.

  - Upload an HDD image and set its properties in one go (ADVANCED):
    ionosctl img upload -i ubuntu.vmdk -l de/fra --rename ubuntu-24.04 --name "Ubuntu 24.04" --licence-type LINUX --cloud-init V1
    Uploads to ftp://ftp-fra.ionos.com/hdd-images/ubuntu-24.04.vmdk, then PATCHes the resulting image with a display name, LINUX licence-type and cloud-init V1.

  - Upload a Confidential Computing (SEV-SNP) image:
    ionosctl img upload -i coco.qcow2 -l de/fra --confidential --name "coco-guest" --licence-type LINUX
    Uploads to the confidential-images/ directory; cloud-init NONE, legacy BIOS off and all hot-plug are forced.

  - Use a custom FTP server:
    ionosctl img upload -i image.iso --ftp-url "ftp://myftp.example" --crt-path certificates/my-server-crt.pem --skip-update

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --application-type string   Application pre-installed on the image (e.g. an MSSQL edition). Only PUBLIC images may set a value other than UNKNOWN, so this is a no-op on uploaded (private) images. Can be one of: MSSQL-2019-Web, MSSQL-2019-Standard, MSSQL-2019-Enterprise, MSSQL-2022-Web, MSSQL-2022-Standard, MSSQL-2022-Enterprise, UNKNOWN (default "UNKNOWN")
      --cloud-init string         Whether servers built from this image accept cloud-init user-data for first-boot provisioning. V1 enables the cloud-init datasource; NONE disables it. Confidential Computing images are forced to NONE. Can be one of: V1, NONE (default "V1")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ImageId Name ImageAliases Location LicenceType ImageType CloudInit CreatedDate Size Description Public CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType RequiredFeatures]
      --confidential              Upload as a Confidential Computing (SEV-SNP) image, to the confidential-images/ directory. Requires a QCOW2 (.qcow/.qcow2) image with an embedded LAUNCH_ARTIFACTS partition. The API only accepts a fixed property set here, so this forces cloud-init NONE, legacy BIOS off, and all hot-plug/hot-unplug off; passing conflicting flags is rejected
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cpu-hot-plug              Guest supports adding CPU cores while running (no reboot) (default true)
      --cpu-hot-unplug            Guest supports removing CPU cores while running. Only valid if CPU hot-plug is also enabled
      --crt-path string           Path to a PEM file with the FTP server's certificate, to trust a self-signed custom server without --skip-verify. Not needed for IONOS FTP servers
  -D, --depth int                 Level of detail for response objects (default 1)
  -d, --description string        Free-text description to set on the uploaded image
      --disc-scsi-hot-plug        Guest supports attaching a SCSI disk while running (no reboot) (default true)
      --disc-scsi-hot-unplug      Guest supports detaching a SCSI disk while running. Not supported on Windows guests
      --disc-virtio-hot-plug      Guest supports attaching a Virt-IO disk while running (no reboot) (default true)
      --disc-virtio-hot-unplug    Guest supports detaching a Virt-IO disk while running. Not supported on Windows guests
      --expose-serial             Expose the attached disk's serial id to the guest. Some OSes/software need it; note it can influence licensed-software (e.g. Windows) behavior
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
      --ftp-port int              TCP port of the FTP server (1-65535, default 21). Only valid together with a custom --ftp-url, for servers on a non-standard port (default 21)
      --ftp-url string            FTP server host. Keep the %s placeholder to have it replaced per --location (default targets IONOS servers ftp-<region>.ionos.com); or give a full custom host to upload elsewhere (default "ftp-%s.ionos.com")
  -h, --help                      Print usage
  -i, --image strings             Comma-separated path(s) to local image file(s) to upload (absolute or relative to cwd). Extension must be one of: .iso .img .vmdk .vhd .vhdx .cow .qcow .qcow2 .raw .vpc .vdi (required)
      --licence-type string       OS licence type. Determines how IONOS bills and configures the guest (e.g. Windows editions are licensed). Use LINUX/RHEL for Linux, WINDOWS2016/2019/2022/2025 for the matching Windows Server, OTHER/UNKNOWN otherwise. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER (default "UNKNOWN")
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location strings          Comma-separated IONOS location(s) to upload each image to; the file is uploaded to every location. One of de/fra, de/fra/2, es/vit, gb/lhr, gb/bhx, fr/par, us/las, us/ewr, us/mci, de/txl, de/fkb. Required unless a custom --ftp-url (without a %s placeholder) is given (required)
  -n, --name string               Human-friendly display name to give the uploaded image (does not have to be unique)
      --nic-hot-plug              Guest supports attaching a NIC while running (no reboot) (default true)
      --nic-hot-unplug            Guest supports detaching a NIC while running. Only valid if NIC hot-plug is also enabled
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
      --ram-hot-plug              Guest supports adding RAM while running (no reboot) (default true)
      --ram-hot-unplug            Guest supports removing RAM while running. Only valid if RAM hot-plug is also enabled
      --rename strings            Comma-separated new names for the uploaded images (positional, one per --image). The file extension is appended automatically; if you include it here (e.g. 'myimage.iso') it is not duplicated. Defaults to each file's base name. Also becomes the image-alias
      --require-legacy-bios       Boot the image in legacy BIOS mode instead of UEFI. Set false for images that require/expect UEFI (default true)
      --skip-update               Only upload the file(s) over FTP; skip the follow-up API step that polls for the image and PATCHes it with the property flags. Use when you just want the bytes on the server
      --skip-verify               Do NOT verify the FTP server's TLS certificate. WARNING: enables man-in-the-middle attacks. Prefer --crt-path for a self-signed custom server
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

