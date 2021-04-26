---
description: Delete an IPBlock
---

# Delete

## Usage

```text
ionosctl ipblock delete [flags]
```

## Description

Use this command to delete a specified IPBlock.

You can wait for the action to be executed using `--wait` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* IPBlock Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings        Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force               Force command to execute without user input
  -h, --help                help for delete
      --ipblock-id string   The unique IPBlock Id [Required flag]
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
      --timeout int         Timeout option for the IPBlock to be deleted [seconds] (default 60)
      --wait                Wait for the IPBlock to be deleted
```

## Examples

```text
ionosctl ipblock delete --ipblock-id bf932826-d71b-4759-a7d0-0028261c1e8d --wait 
Warning: Are you sure you want to delete ipblock (y/N) ? 
y
Waiting for request: 6b1aa258-799f-4712-9f90-ba4494d84026
RequestId: 6b1aa258-799f-4712-9f90-ba4494d84026
Status: Command ipblock delete and request have been successfully executed
```

