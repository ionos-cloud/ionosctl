---
description: List attached Volumes from a Server
---

# ListVolumes

## Usage

```text
ionosctl server list-volumes [flags]
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
      --datacenter-id string   The unique Data Center Id [Required flag]
  -h, --help                   help for list-volumes
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id [Required flag]
```

## Examples

```text
ionosctl server list-volumes --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --server-id 2bf04e0d-86e4-4f13-b405-442363b25e28 
VolumeId                               Name   Size   Type   LicenceType   State       Image
1ceb4b02-ed41-4651-a90b-9a30bc284e74   test   10GB   HDD    LINUX         AVAILABLE
```

