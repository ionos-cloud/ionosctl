---
description: Delete a LAN
---

# LanDelete

## Usage

```text
ionosctl lan delete [flags]
```

## Aliases

For `lan` command:

```text
[l]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified LAN from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LanId Name Public PccId State] (default [LanId,Name,Public,PccId,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for delete
  -i, --lan-id string          The unique LAN Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for LAN deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for Request for LAN deletion to be executed
```

## Examples

```text
ionosctl lan delete --datacenter-id DATACENTER_ID --lan-id LAN_ID

ionosctl lan delete --datacenter-id DATACENTER_ID --lan-id LAN_ID --wait-for-request
```

