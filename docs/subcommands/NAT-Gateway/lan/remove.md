---
description: "Remove a NAT Gateway Lan"
---

# NatgatewayLanRemove

## Usage

```text
ionosctl natgateway lan remove [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a specified NAT Gateway Lan from a NAT Gateway.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Lan Id

## Options

```text
  -a, --all                    Remove all NAT Gateway Lans.
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayLanId GatewayIps] (default [NatGatewayLanId,GatewayIps])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -i, --lan-id int             The unique LAN Id (required) (default 1)
      --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout in seconds for polling the request (default 60)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait                   Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl natgateway lan remove --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id LAN_ID
```

