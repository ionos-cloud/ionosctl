---
description: Delete an IpBlock
---

# IpblockDelete

## Usage

```text
ionosctl ipblock delete [flags]
```

## Aliases

For `ipblock` command:
```text
[block ipb]
```

## Description

Use this command to delete a specified IpBlock.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* IpBlock Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [IpBlockId Name Location Size Ips State] (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                help for delete
  -i, --ipblock-id string   The unique IpBlock Id (required)
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout option for Request for IpBlock deletion [seconds] (default 60)
  -w, --wait-for-request    Wait for the Request for IpBlock deletion to be executed
```

## Examples

```text
ionosctl ipblock delete --ipblock-id bf932826-d71b-4759-a7d0-0028261c1e8d --wait-for-request 
Warning: Are you sure you want to delete ipblock (y/N) ? 
y
1.2s Waiting for request... DONE
RequestId: 6b1aa258-799f-4712-9f90-ba4494d84026
Status: Command ipblock delete & wait have been successfully executed
```

