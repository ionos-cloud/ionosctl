---
description: "Update a private image's properties and hardware capabilities"
---

# ImageUpdate

## Usage

```text
ionosctl compute image update [flags]
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

Update the metadata and advertised hardware capabilities of one of your PRIVATE images (PUBLIC IONOS images cannot be modified). Typical uses are correcting the licence-type after an FTP upload, toggling cloud-init support, or declaring which resources can be hot-plugged.

The hot-plug / hot-unplug flags describe what a server booted from this image will support at runtime, i.e. whether the guest OS can have that resource added (hot-plug) or removed (hot-unplug) while the VM is running, without a reboot. Set them to match what your OS and drivers actually support; enabling a capability the guest cannot handle leads to failed or ignored operations.

Constraints:
  * You can only enable hot-UNPLUG for a resource whose hot-PLUG you also enabled.
  * Disk hot-unplug (--disc-virtio-hot-unplug, --disc-scsi-hot-unplug) is not supported for Windows guests.
  * --application-type only accepts a value other than UNKNOWN on PUBLIC images, so it is effectively a no-op on your private images.

Note: all boolean capability flags carry a default, so every one is sent on update; pass the flags explicitly to set the values you want.

Required values to run command:

* Image Id

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --application-type string   Application pre-installed on the image (e.g. an MSSQL edition). Only PUBLIC images may set a value other than UNKNOWN, so this is a no-op on your private images. Can be one of: MSSQL-2019-Web, MSSQL-2019-Standard, MSSQL-2019-Enterprise, MSSQL-2022-Web, MSSQL-2022-Standard, MSSQL-2022-Enterprise, UNKNOWN (default "UNKNOWN")
      --cloud-init string         Whether servers built from this image accept cloud-init user-data for first-boot provisioning. V1 enables the cloud-init datasource; NONE disables it. Can be one of: V1, NONE (default "V1")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ImageId Name ImageAliases Location LicenceType ImageType CloudInit CreatedDate Size Description Public CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType RequiredFeatures]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cpu-hot-plug              Guest supports adding CPU cores while running (no reboot) (default true)
      --cpu-hot-unplug            Guest supports removing CPU cores while running. Only valid if CPU hot-plug is also enabled
  -D, --depth int                 Level of detail for response objects (default 1)
  -d, --description string        Free-text description of the image
      --disc-scsi-hot-plug        Guest supports attaching a SCSI disk while running (no reboot) (default true)
      --disc-scsi-hot-unplug      Guest supports detaching a SCSI disk while running. Not supported on Windows guests
      --disc-virtio-hot-plug      Guest supports attaching a Virt-IO disk while running (no reboot) (default true)
      --disc-virtio-hot-unplug    Guest supports detaching a Virt-IO disk while running. Not supported on Windows guests
      --expose-serial             Expose the attached disk's serial id to the guest. Some OSes/software need it; note it can influence licensed-software (e.g. Windows) behavior
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
  -i, --image-id string           The unique Image Id (required)
      --licence-type string       OS licence type. Determines how IONOS bills and configures the guest (e.g. Windows editions are licensed). Use LINUX/RHEL for Linux, WINDOWS2016/2019/2022/2025 for the matching Windows Server, OTHER/UNKNOWN otherwise. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER (default "UNKNOWN")
      --limit int                 Maximum number of items to return per request (default 50)
  -n, --name string               Human-friendly display name of the image (does not have to be unique)
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
      --require-legacy-bios       Boot the image in legacy BIOS mode instead of UEFI. Set false for images that require/expect UEFI (default true)
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Correct the licence type and name of an uploaded image
ionosctl compute image update --image-id IMAGE_ID --name "ubuntu-24.04-custom" --licence-type LINUX

# Declare a Linux image that supports CPU/RAM/NIC hot-plug but no disk hot-unplug, with cloud-init
ionosctl compute image update --image-id IMAGE_ID --licence-type LINUX --cloud-init V1 \
  --cpu-hot-plug=true --ram-hot-plug=true --nic-hot-plug=true \
  --disc-virtio-hot-unplug=false --disc-scsi-hot-unplug=false
```

