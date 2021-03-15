---
description: Delete a LAN
---

# Delete

## Usage

```text
ionosctl lan delete [flags]
```

## Description

Use this command to delete a specified LAN from a Virtual Data Center.

You can wait for the action to be executed using `--wait` option. You can force the command to execute without user input using `--ignore-stdin` option.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [LanId,Name,Public])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for delete
      --ignore-stdin           Force command to execute without user input
      --lan-id string          The unique LAN Id [Required flag]
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --timeout int            Timeout option [seconds] (default 60)
      --wait                   Wait for LAN to be deleted
```

## Examples

```text
ionosctl lan delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 4
Warning: Are you sure you want to delete lan (y/N) ? y
RequestId: bd5ffcf4-1b05-4cb2-917b-a0140d5f7a2b
Status: Command lan delete has been successfully executed

ionosctl lan delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 3 --wait 
Warning: Are you sure you want to delete lan (y/N) ? y
Waiting for request: e65fc2fe-8005-48a5-9d06-f1a4f8bc9ef1
RequestId: e65fc2fe-8005-48a5-9d06-f1a4f8bc9ef1
Status: Command lan delete and request have been successfully executed
```

