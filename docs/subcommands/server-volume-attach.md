---
description: Attach a Volume to a Server
---

# ServerVolumeAttach

## Usage

```text
ionosctl server volume attach [flags]
```

## Description

Use this command to attach a pre-existing Volume to a Server.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for attach
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Volume attachment [seconds] (default 60)
      --volume-id string       The unique Volume Id (required)
  -w, --wait-for-request       Wait for the Request for Volume attachment to be executed
```

## Examples

```text
ionosctl server volume attach --datacenter-id aa8e07a2-287a-4b45-b5e9-94761750a53c --server-id 1dc7c6a8-5ab3-4fa8-83e7-9d989bd52ffa --volume-id 101291d1-2227-432a-9773-97b5ace7b8d3 
VolumeId                               Name   Size   Type   LicenceType   State   Image
101291d1-2227-432a-9773-97b5ace7b8d3   test   10GB   HDD    LINUX         BUSY    
RequestId: e8ad392c-006a-487a-8852-c38b6e7f7ad7
Status: Command volume attach has been successfully executed
```

