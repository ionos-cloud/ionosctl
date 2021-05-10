---
description: Update an IpBlock
---

# IpblockUpdate

## Usage

```text
ionosctl ipblock update [flags]
```

## Description

Use this command to update the properties of an existing IpBlock.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* IpBlock Id

## Options

```text
  -u, --api-url string        Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings          Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                 Force command to execute without user input
  -h, --help                  help for update
      --ipblock-id string     The unique IpBlock Id (required)
      --ipblock-name string   Name of the IpBlock
  -o, --output string         Desired output format [text|json] (default "text")
  -q, --quiet                 Quiet output
      --timeout int           Timeout option for Request for IpBlock update [seconds] (default 60)
      --wait-for-request      Wait for the Request for IpBlock update to be executed
```

## Examples

```text
ionosctl ipblock update --ipblock-id bf932826-d71b-4759-a7d0-0028261c1e8d --ipblock-name demo
IpBlockId                              Name   Location   Size   Ips         State
bf932826-d71b-4759-a7d0-0028261c1e8d   demo   us/las     1      [x.x.x.x]   BUSY
RequestId: 5864afe5-4df5-4843-b548-4489857dc3c5
Status: Command ipblock update has been successfully executed
```

