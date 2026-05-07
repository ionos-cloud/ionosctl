---
description: "Update a Application Load Balancer Forwarding Rule"
---

# ApplicationloadbalancerRuleUpdate

## Usage

```text
ionosctl compute applicationloadbalancer rule update [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `rule` command:

```text
[r forwardingrule]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified Application Load Balancer Forwarding Rule from a Application Load Balancer. You can also update Health Check settings.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id

## Options

```text
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --client-timeout int                  The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds). (default 50)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ForwardingRuleId Name Protocol ListenerIp ListenerPort ClientTimeout ServerCertificates State]
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int                           Level of detail for response objects (default 1)
  -F, --filters strings                     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --limit int                           Maximum number of items to return per request (default 50)
      --listener-ip ip                      Listening (inbound) IP.
      --listener-port int                   Listening (inbound) port number; valid range is 1 to 65535. (default 8080)
  -n, --name string                         The name of the Application Load Balancer forwarding rule.
      --no-headers                          Don't print table headers when table output is used
      --offset int                          Number of items to skip before starting to collect the results
      --order-by string                     Property to order the results by
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
      --query string                        JMESPath query string to filter the output
  -q, --quiet                               Quiet output
  -i, --rule-id string                      The unique ForwardingRule Id (required)
      --server-certificates strings         Server Certificates
  -t, --timeout int                         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                                Wait for the resource to reach AVAILABLE state after the command completes
```

## Examples

```text
ionosctl compute alb rule update --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FORWARDINGRULE_ID --name NAME
```

