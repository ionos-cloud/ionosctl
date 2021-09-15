---
description: Create a Application Load Balancer Forwarding Rule
---

# ApplicationloadbalancerRuleCreate

## Usage

```text
ionosctl applicationloadbalancer rule create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create a Application Load Balancer Forwarding Rule in a specified Application Load Balancer. You can also set Health Check Settings for Forwarding Rule.

You can wait for the Request to be executed using `--wait-for-request` or `-w` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Listener Ip
* Listener Port

## Options

```text
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ForwardingRuleId Name Protocol ListenerIp ListenerPort ServerCertificates State] (default [ForwardingRuleId,Name,Protocol,ListenerIp,ListenerPort,ServerCertificates,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id (required)
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --listener-ip string                  Listening IP (required)
      --listener-port int                   Listening port number. Range: 1 to 65535 (required) (default 8080)
  -n, --name string                         The name for the Forwarding Rule (default "Unnamed Forwarding Rule")
  -o, --output string                       Desired output format [text|json] (default "text")
  -p, --protocol string                     Protocol of the balancing (default "HTTP")
  -q, --quiet                               Quiet output
      --server-certificates strings         Server Certificates
  -t, --timeout int                         Timeout option for Request for Forwarding Rule creation [seconds] (default 300)
  -v, --verbose                             Print step-by-step process when running command
  -w, --wait-for-request                    Wait for the Request for Forwarding Rule creation to be executed
```

## Examples

```text
ionosctl applicationloadbalancer rule create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --listener-ip LISTENER_IP --listener-port LISTENER_PORT
```

