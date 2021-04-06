---
description: Update an IPBlock
---

# Update

## Usage

```text
ionosctl ipblock update [flags]
```

## Description

Use this command to update a specified IPBlock.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* IPBlock Id

## Options

```text
  -u, --api-url string        Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings          Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                  help for update
      --ignore-stdin          Force command to execute without user input
      --ipblock-id string     The unique IPBlock Id [Required flag]
      --ipblock-name string   Name of the IPBlock
  -o, --output string         Desired output format [text|json] (default "text")
  -q, --quiet                 Quiet output
      --timeout int           Timeout option for the IPBlock to be updated [seconds] (default 60)
      --wait                  Wait for the IPBlock to be updated
```

## Examples

```text
ionosctl ipblock update --ipblock-id bf932826-d71b-4759-a7d0-0028261c1e8d --ipblock-name demo
IpBlockId                              Name   Location   Size   Ips         State
bf932826-d71b-4759-a7d0-0028261c1e8d   demo   us/las     1      [x.x.x.x]   BUSY
RequestId: 5864afe5-4df5-4843-b548-4489857dc3c5
Status: Command ipblock update has been successfully executed
```

