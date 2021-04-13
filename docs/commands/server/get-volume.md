---
description: Get an attached Volume from a Server
---

# GetVolume

## Usage

```text
ionosctl server get-volume [flags]
```

## Description

Use this command to retrieve information about an attached Volume.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
  -h, --help                   help for get-volume
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id [Required flag]
      --volume-id string       The unique Volume Id [Required flag]
```

## Examples

```text
ionosctl volume attach get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --volume-id 15546173-a100-4851-8bc4-872ec6bbee32 
VolumeId                               Name         Size   Type   LicenseType   State       Image
15546173-a100-4851-8bc4-872ec6bbee32   demoVolume   10GB   HDD    LINUX         AVAILABLE
```

