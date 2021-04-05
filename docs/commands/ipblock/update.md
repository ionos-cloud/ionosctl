---
description: Update an IpBlock
---

# Update

## Usage

```text
ionosctl ipblock update [flags]
```

## Description

Use this command to update a specified IpBlock.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* IpBlock Id

## Options

```text
  -u, --api-url string        Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings          Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                  help for update
      --ignore-stdin          Force command to execute without user input
      --ipblock-id string     The unique IpBlock Id [Required flag]
      --ipblock-ips strings   Ips of the IpBlock
      --ipblock-name string   Name of the IpBlock
  -o, --output string         Desired output format [text|json] (default "text")
  -q, --quiet                 Quiet output
      --timeout int           Timeout option for the IpBlock to be updated [seconds] (default 60)
      --wait                  Wait for the IpBlock to be updated
```

