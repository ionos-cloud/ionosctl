---
description: Delete a LAN
---

# LanDelete

## Usage

```text
ionosctl lan delete [flags]
```

## Description

Use this command to delete a specified LAN from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [LanId,Name,Public,PccId])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --force                  Force command to execute without user input
  -h, --help                   help for delete
      --lan-id string          The unique LAN Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --timeout int            Timeout option for Request for LAN deletion [seconds] (default 60)
      --wait-for-request       Wait for Request for LAN deletion to be executed
```

## Examples

```text
ionosctl lan delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 4
Warning: Are you sure you want to delete lan (y/N) ? y
RequestId: bd5ffcf4-1b05-4cb2-917b-a0140d5f7a2b
Status: Command lan delete has been successfully executed

ionosctl lan delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 3 --wait-for-request 
Warning: Are you sure you want to delete lan (y/N) ? y
1.2s Waiting for request... DONE
RequestId: e65fc2fe-8005-48a5-9d06-f1a4f8bc9ef1
Status: Command lan delete & wait have been successfully executed
```

