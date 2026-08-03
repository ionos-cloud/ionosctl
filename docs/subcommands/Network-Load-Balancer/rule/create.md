---
description: "Create a Network Load Balancer Forwarding Rule"
---

# NetworkloadbalancerRuleCreate

## Usage

```text
ionosctl compute networkloadbalancer rule create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create a forwarding rule on a Network Load Balancer. The rule listens on --listener-ip (which must be one of the NLB's own IPs) and --listener-port, then balances accepted TCP connections across the targets you later add with `nlb rule target add`.

Pick a balancing algorithm with --algorithm:
  ROUND_ROBIN (default), LEAST_CONNECTION, RANDOM, SOURCE_IP (client-IP affinity).

The health-check flags tune resilience: --connect-timeout bounds how long a new connection to a target may take, --client-timeout / --target-timeout bound client- and target-side inactivity, and --retries sets how many times to retry a target after a connection failure before considering it down.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Listener Ip
* Listener Port

## Options

```text
      --algorithm string                Balancing algorithm used to pick a target per connection: ROUND_ROBIN, LEAST_CONNECTION, RANDOM, or SOURCE_IP (client-IP affinity) (default "ROUND_ROBIN")
  -u, --api-url string                  Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --client-timeout int              [Health Check] Maximum client-side inactivity, in milliseconds, before the connection is closed (client expected to acknowledge or send data) (default 5000)
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [ForwardingRuleId Name Algorithm Protocol ListenerIp ListenerPort State ClientTimeout ConnectTimeout TargetTimeout Retries]
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --connection-timeout int          [Health Check] Maximum time, in milliseconds, to wait for a connection to a target VM to succeed (default 5000)
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int                       Level of detail for response objects (default 1)
  -F, --filters strings                 Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --limit int                       Maximum number of items to return per request (default 50)
      --listener-ip ip                  Inbound IP the rule listens on. Must be one of the NLB's own IPs (--ips) (required)
      --listener-port string            Inbound TCP port the rule listens on. Range: 1 to 65535 (required)
  -n, --name string                     The name for the Forwarding Rule (default "Unnamed Forwarding Rule")
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      Don't print table headers when table output is used
      --offset int                      Number of items to skip before starting to collect the results
      --order-by string                 Property to order the results by
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
      --query string                    JMESPath query string to filter the output
  -q, --quiet                           Quiet output
      --retries int                     [Health Check] Number of times to retry a target after a connection failure before marking it down. Range: 0 to 65535 (default 3)
      --target-timeout int              [Health Check] Maximum target-side inactivity, in milliseconds, before the connection is closed (default 5000)
  -t, --timeout int                     Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                   Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                            Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a rule that listens on port 80 (defaults to ROUND_ROBIN)
ionosctl compute networkloadbalancer rule create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --listener-ip 203.0.113.10 --listener-port 80

# Create a rule with client-IP affinity and custom health-check timeouts
ionosctl compute networkloadbalancer rule create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --name "https" --listener-ip 203.0.113.10 --listener-port 443 --algorithm SOURCE_IP --connect-timeout 5000 --retries 5
```

