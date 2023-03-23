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
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --client-timeout int                  The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds). (default 50)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [ForwardingRuleId Name Protocol ListenerIp ListenerPort ServerCertificates State] (default [ForwardingRuleId,Name,Protocol,ListenerIp,ListenerPort,ServerCertificates,State])
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
      --listener-ip ip                      Listening (inbound) IP. It must be assigned to the listener NIC of Application Load Balancer. (required)
      --listener-port int                   Listening (inbound) port number; valid range is 1 to 65535. (required) (default 8080)
  -n, --name string                         The name of the Application Load Balancer forwarding rule. (default "Unnamed Forwarding Rule")
  -p, --protocol string                     Balancing protocol. (default "HTTP")
      --server-certificates strings         Server Certificates
  -t, --timeout int                         Timeout option for Request for Forwarding Rule creation [seconds] (default 300)
  -w, --wait-for-request                    Wait for the Request for Forwarding Rule creation to be executed
```

## Examples

```text
ionosctl applicationloadbalancer rule create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --listener-ip LISTENER_IP --listener-port LISTENER_PORT
```

