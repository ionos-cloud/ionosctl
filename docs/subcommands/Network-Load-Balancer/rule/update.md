---
description: "Update a Network Load Balancer Forwarding Rule"
---

# NetworkloadbalancerRuleUpdate

## Usage

```text
ionosctl networkloadbalancer rule update [flags]
```

## Aliases

For `networkloadbalancer` command:

```text
[nlb]
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

Use this command to update a specified Network Load Balancer Forwarding Rule from a Network Load Balancer. You can also update Health Check settings.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id

## Options

```text
      --algorithm string                Algorithm for the balancing (default "ROUND_ROBIN")
  -u, --api-url string                  Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --client-timeout int              [Health Check] ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data (default 5000)
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [ForwardingRuleId Name Algorithm Protocol ListenerIp ListenerPort State ClientTimeout ConnectTimeout TargetTimeout Retries] (default [ForwardingRuleId,Name,Algorithm,Protocol,ListenerIp,ListenerPort,State])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --connection-timeout int          [Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed (default 5000)
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int32                     Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --listener-ip ip                  Listening IP (required)
      --listener-port string            Listening port number. Range: 1 to 65535 (required)
  -n, --name string                     The name for the Forwarding Rule
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      Don't print table headers when table output is used
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
  -q, --quiet                           Quiet output
      --retries int                     [Health Check] Retries specifies the number of retries to perform on a target VM after a connection failure. Range: 0 to 65535 (default 3)
  -i, --rule-id string                  The unique ForwardingRule Id (required)
      --target-timeout int              [Health Check] TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side (default 5000)
  -t, --timeout int                     Timeout option for Request for Forwarding Rule update [seconds] (default 300)
  -v, --verbose                         Print step-by-step process when running command
  -w, --wait-for-request                Wait for the Request for Forwarding Rule update to be executed
```

## Examples

```text
ionosctl nlb rule update --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FORWARDINGRULE_ID --name NAME
```

