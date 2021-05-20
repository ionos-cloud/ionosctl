---
description: Create a LAN
---

# LanCreate

## Usage

```text
ionosctl lan create [flags]
```

## Description

Use this command to create a new LAN within a Virtual Data Center on your account. The name, the public option and the Private Cross-Connect Id can be set.

NOTE: IP Failover is configured after LAN creation using an update command.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [LanId,Name,Public,PccId,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for create
  -n, --name string            The name of the LAN
  -o, --output string          Desired output format [text|json] (default "text")
      --pcc-id string          The unique Id of the Private Cross-Connect the LAN will connect to
      --public                 Indicates if the LAN faces the public Internet (true) or not (false)
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for LAN creation [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for LAN creation to be executed
```

## Examples

```text
ionosctl lan create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-name demoLan
LanId   Name      Public   PccId
4       demoLan   false
RequestId: da824a69-a12a-4153-b302-a797b3581c2b
Status: Command lan create has been successfully executed
```

