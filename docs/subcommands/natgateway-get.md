---
description: Get a NAT Gateway
---

# NatgatewayGet

## Usage

```text
ionosctl natgateway get [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified NAT Gateway from a Virtual Data Center. You can also wait for NAT Gateway to get in AVAILABLE state using `--wait-for-state` option.

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
  -h, --help                   help for get
  -i, --natgateway-id string   The unique NatGateway Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for waiting for NAT Gateway to be in AVAILABLE state [seconds] (default 60)
  -v, --verbose                see step by step process when running a command
  -W, --wait-for-state         Wait for specified NAT Gateway to be in AVAILABLE state
```

## Examples

```text
ionosctl natgateway get --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID
```

