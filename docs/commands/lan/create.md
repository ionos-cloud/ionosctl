---
description: Create a LAN
---

# Create

## Usage

```text
ionosctl lan create [flags]
```

## Description

Use this command to create a new LAN within a Virtual Data Center on your account. The name and public option can be set. Please Note: IP Failover is configured after LAN creation using an update command.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [LanId,Name,Public])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for create
      --ignore-stdin           Force command to execute without user input
      --lan-name string        The name of the LAN
      --lan-public             Indicates if the LAN faces the public Internet (true) or not (false)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --timeout int            Timeout option [seconds] (default 60)
  -v, --verbose                Enable verbose output
      --wait                   Wait for LAN to be created
```

## Examples

```text
ionosctl lan create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-name demoLan
LanId   Name      Public
4       demoLan   false
RequestId: da824a69-a12a-4153-b302-a797b3581c2b
Status: Command lan create has been successfully executed
```

