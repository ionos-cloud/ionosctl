---
description: "Delete a LAN"
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
  -a, --all                    Delete all Lans from a Virtual Data Center.
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LanId Name Public PccId IPv6CidrBlock State DatacenterId] (default [LanId,Name,Public,PccId,IPv6CidrBlock,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -i, --lan-id string          The unique LAN Id (required)
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for LAN deletion [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait                   Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl lan delete --datacenter-id DATACENTER_ID --lan-id LAN_ID

ionosctl lan delete --datacenter-id DATACENTER_ID --lan-id LAN_ID --wait-for-request
```

