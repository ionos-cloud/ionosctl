---
description: "Get a NAT Gateway"
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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayId Name PublicIps State] (default [NatGatewayId,Name,PublicIps,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -i, --natgateway-id string   The unique NatGateway Id (required)
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout duration       Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait                   Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl natgateway get --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID
```

