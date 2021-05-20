---
description: Delete a Server
---

# ServerDelete

## Usage

```text
ionosctl server delete [flags]
```

## Aliases

For `server` command:
```text
[s]
```

## Description

Use this command to delete a specified Server from a Virtual Data Center.

NOTE: This will not automatically remove the storage Volume(s) attached to a Server.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -F, --format strings         Collection of fields to be printed on output (default [ServerId,Name,AvailabilityZone,Cores,Ram,CpuFamily,VmState,State])
  -h, --help                   help for delete
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Server deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for Server deletion to be executed
```

## Examples

```text
ionosctl server delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id f45f435e-8d6c-4170-ab90-858b59dab9ff 
Warning: Are you sure you want to delete server (y/N) ? Y
RequestId: 1f00c6d9-072a-4dd0-8c09-c46f2f20a230
Status: Command server delete has been successfully executed

ionosctl server delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 35201d04-0ea2-43e7-abc4-56f92737bb9d --force 
RequestId: f596caba-78b7-4c99-8c9d-56198d3754b6
Status: Command server delete has been successfully executed
```

