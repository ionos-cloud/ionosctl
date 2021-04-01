---
description: Detach a Volume from a Server
---

# Detach

## Usage

```text
ionosctl volume detach [flags]
```

## Description

Use this command to detach a Volume from a Server.

You can wait for the action to be executed using `--wait` option. You can force the command to execute without user input using `--ignore-stdin` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for detach
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id [Required flag]
      --timeout int            Timeout option for Server to be detached from a Server [seconds] (default 60)
      --volume-id string       The unique Volume Id [Required flag]
      --wait                   Wait for Volume to detach from Server
```

## Examples

```text
ionosctl volume detach --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --volume-id 15546173-a100-4851-8bc4-872ec6bbee32 
Warning: Are you sure you want to detach volume (y/N) ? y
RequestId: bb4d79ef-129d-4e39-8f5c-519b7cefbc54
Status: Command volume detach has been successfully executed
```

