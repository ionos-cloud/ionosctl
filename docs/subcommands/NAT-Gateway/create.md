---
description: "Create a NAT Gateway"
---

# NatgatewayCreate

## Usage

```text
ionosctl natgateway create [flags]
```

## Aliases

For `natgateway` command:

```text
[nat ng]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a NAT Gateway in a specified Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* IPs

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
      --ips strings            Collection of public reserved IP addresses of the NAT Gateway (required)
  -n, --name string            Name of the NAT Gateway (default "NAT Gateway")
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout duration       Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait                   Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl natgateway create --datacenter-id DATACENTER_ID --name NAME --ips IP_1,IP_2
```

