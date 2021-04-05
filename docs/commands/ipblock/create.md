---
description: Create/Reserve an IpBlock
---

# Create

## Usage

```text
ionosctl ipblock create [flags]
```

## Description

Use this command to create/reserve an IpBlock in a specified location. The name, size, location options can be set.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* IpBlock Location
* IpBlock Name
* IpBlock Size

## Options

```text
  -u, --api-url string            Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings              Columns to be printed in the standard output (default [IpBlockId,Name,Location,Size,Ips,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                      help for create
      --ignore-stdin              Force command to execute without user input
      --ipblock-location string   Location of the IpBlock [Required flag]
      --ipblock-name string       Name of the IpBlock [Required flag]
      --ipblock-size int          Size of the IpBlock [Required flag] (default 2)
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
      --timeout int               Timeout option for the IpBlock to be created [seconds] (default 60)
      --wait                      Wait for the IpBlock to be created
```

