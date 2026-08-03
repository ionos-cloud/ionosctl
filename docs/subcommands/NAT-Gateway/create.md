---
description: "Create a NAT Gateway"
---

# NatgatewayCreate

## Usage

```text
ionosctl compute natgateway create [flags]
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

Use this command to create a NAT Gateway in a specified Virtual Data Center. The gateway provides source-NAT (SNAT) outbound internet access for servers on private LANs, masquerading their traffic behind the gateway's public IPs.

Creating the gateway only allocates it and its public IPs. To make traffic actually flow you still need to attach it to a private LAN (`natgateway lan add`) and add at least one SNAT rule (`natgateway rule create`).

The addresses passed to `--ips` must be public IPs you have already reserved in the same location as the datacenter (see `ionosctl compute ipblock`); arbitrary or in-use addresses are rejected.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state before returning.

Required values to run command:

* Data Center Id
* IPs

## Options

```text
  -u, --api-url string                 Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings                   Set of columns to be printed on output 
                                       Available columns: [NatGatewayId Name PublicIps State DatacenterId]
  -c, --config string                  Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string           The unique Data Center Id (required)
  -D, --depth int                      Level of detail for response objects (default 1)
  -F, --filters strings                Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                          Force command to execute without user input
  -h, --help                           Print usage
      --ips ionosctl compute ipblock   Comma-separated public IP addresses the gateway masquerades outbound traffic behind. Must be IPs you have already reserved in the same location as the datacenter (see ionosctl compute ipblock). SNAT rules on this gateway can only reference IPs listed here (required)
      --limit int                      Maximum number of items to return per request (default 50)
  -n, --name string                    Human-friendly name for the NAT Gateway (default "NAT Gateway")
      --no-headers                     Don't print table headers when table output is used
      --offset int                     Number of items to skip before starting to collect the results
      --order-by string                Property to order the results by
  -o, --output string                  Desired output format [text|json|api-json] (default "text")
      --query string                   JMESPath query string to filter the output
  -q, --quiet                          Quiet output
  -t, --timeout int                    Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                  Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                           Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a NAT Gateway with a single reserved public IP
ionosctl compute natgateway create --datacenter-id DATACENTER_ID --name my-gateway --ips 203.0.113.10

# Create a NAT Gateway with two public IPs and wait until it is AVAILABLE
ionosctl compute natgateway create --datacenter-id DATACENTER_ID --name my-gateway --ips 203.0.113.10,203.0.113.11 --wait
```

