---
description: Stop a Server
---

# Ionosctl Server Stop

## Usage

```text
ionosctl server stop [flags]
```

## Description

Use this command to stop specified Server from a Data Center.

You can wait for the action to be executed using `--wait` option. You can force the command to execute without user input using `--ignore-stdin` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [DatacenterId,Name,Location])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for stop
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id [Required flag]
      --timeout int            Timeout option [seconds] (default 60)
      --wait                   Wait for Server to stop
```

## Examples

```text
ionosctl server stop --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa
Warning: Are you sure you want to stop server (y/N) ? y
RequestId: 8c06523d-8838-4409-aee3-68c042f5a256
Status: Command server stop has been successfully executed
```

