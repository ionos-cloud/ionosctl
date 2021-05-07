---
description: Start a Server
---

# ServerStart

## Usage

```text
ionosctl server start [flags]
```

## Description

Use this command to start a Server from a Virtual Data Center. If the Server's public IP was deallocated then a new IP will be assigned.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [ServerId,Name,AvailabilityZone,State,Cores,Ram,CpuFamily])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --force                  Force command to execute without user input
  -h, --help                   help for start
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
      --timeout int            Timeout option for Request for Server start [seconds] (default 60)
      --wait-for-request       Wait for the Request for Server start to be executed
```

## Examples

```text
ionosctl server start --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa
Warning: Are you sure you want to start server (y/N) ? y
RequestId: 9f03a764-5f6c-4740-87e2-d9e9589265dc
Status: Command server start has been successfully executed
```

