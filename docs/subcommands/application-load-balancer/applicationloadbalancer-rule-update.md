---
description: Update a Application Load Balancer Forwarding Rule
---

# ApplicationloadbalancerRuleUpdate

## Usage

```text
ionosctl applicationloadbalancer rule update [flags]
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

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id

## Options

```text
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --client-timeout int                  The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds). (default 50)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ForwardingRuleId Name Protocol ListenerIp ListenerPort ServerCertificates State] (default [ForwardingRuleId,Name,Protocol,ListenerIp,ListenerPort,ServerCertificates,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int                           Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --listener-ip string                  Listening (inbound) IP.
      --listener-port int                   Listening (inbound) port number; valid range is 1 to 65535. (default 8080)
  -n, --name string                         The name of the Application Load Balancer forwarding rule.
  -o, --output string                       Desired output format [text|json] (default "text")
  -q, --quiet                               Quiet output
  -i, --rule-id string                      The unique ForwardingRule Id (required)
      --server-certificates strings         Server Certificates
  -t, --timeout int                         Timeout option for Request for Forwarding Rule update [seconds] (default 300)
  -v, --verbose                             Print step-by-step process when running command
  -w, --wait-for-request                    Wait for the Request for Forwarding Rule update to be executed
```

## Examples

```text
ionosctl alb rule update --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FORWARDINGRULE_ID --name NAME
```

