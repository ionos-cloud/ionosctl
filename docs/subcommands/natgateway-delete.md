---
description: Delete a NAT Gateway
---

# NatgatewayDelete

## Usage

```text
ionosctl natgateway delete [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified NAT Gateway from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayId Name PublicIps State] (default [NatGatewayId,Name,PublicIps,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for delete
  -i, --natgateway-id string   The unique NatGateway Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for NAT Gateway deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NAT Gateway deletion to be executed
```

## Examples

```text
ionosctl natgateway delete --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID
```

