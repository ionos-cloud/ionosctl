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
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NatGatewayId Name PublicIps State] (default [NatGatewayId,Name,PublicIps,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            Collection of public reserved IP addresses of the NAT Gateway (required)
      --limit int              Pagination limit: Maximum number of items to return per request (default 50)
  -n, --name string            Name of the NAT Gateway (default "NAT Gateway")
      --no-headers             Don't print table headers when table output is used
      --offset int             Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for NAT Gateway creation [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request       Wait for the Request for NAT Gateway creation to be executed
```

## Examples

```text
ionosctl natgateway create --datacenter-id DATACENTER_ID --name NAME --ips IP_1,IP_2
```

