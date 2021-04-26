---
description: List Volumes
---

# List

## Usage

```text
ionosctl volume list [flags]
```

## Description

Use this command to list all Volumes from a Data Center on your account.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
      --force                  Force command to execute without user input
  -h, --help                   help for list
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl volume list --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d 
VolumeId                               Name         Size   Type   LicenceType   State       Image
ce510144-9bc6-4115-bd3d-b9cd232dd422   demoVolume   20GB   HDD    LINUX         AVAILABLE
```

