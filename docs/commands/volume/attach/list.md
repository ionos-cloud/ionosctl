---
description: List attached Volumes from a Server
---

# Ionosctl Volume Attach List

## Usage

```text
ionosctl volume attach list [flags]
```

## Description

Use this command to get a list of attached Volumes to a Server from a Data Center.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for list
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id [Required flag]
```

## Examples

```text
ionosctl volume attach list --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa 
VolumeId                               Name         Size   Type   LicenseType   State       Image
15546173-a100-4851-8bc4-872ec6bbee32   demoVolume   10GB   HDD    LINUX         AVAILABLE
```

