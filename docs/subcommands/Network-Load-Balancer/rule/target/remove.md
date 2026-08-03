---
description: "Remove a Target from a Network Load Balancer Forwarding Rule"
---

# NetworkloadbalancerRuleTargetRemove

## Usage

```text
ionosctl compute networkloadbalancer rule target remove [flags]
```

## Aliases

For `rule` command:

```text
[r forwardingrule]
```

For `target` command:

```text
[t]
```

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a backend target from a forwarding rule, identified by its --ip and --port. The rule stops forwarding new connections to it immediately; the VM itself is not affected. To temporarily stop traffic without removing a target, prefer --maintenance on the target instead.

Use `--wait` (`-w`) to wait for the removal to complete. You can force the command to execute without user input using `--force` option. Use `--all` to remove every target from the rule.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id
* Target Ip
* Target Port

## Options

```text
  -a, --all                             Remove all targets from the forwarding rule
  -u, --api-url string                  Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [TargetIp TargetPort Weight Check CheckInterval Maintenance]
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int                       Level of detail for response objects (default 1)
  -F, --filters strings                 Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --ip ip                           IP of the backend target VM to remove (required)
      --limit int                       Maximum number of items to return per request (default 50)
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      Don't print table headers when table output is used
      --offset int                      Number of items to skip before starting to collect the results
      --order-by string                 Property to order the results by
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
  -P, --port string                     Port of the backend target to remove. Range: 1 to 65535 (required)
      --query string                    JMESPath query string to filter the output
  -q, --quiet                           Quiet output
      --rule-id string                  The unique ForwardingRule Id (required)
  -t, --timeout int                     Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                   Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                            Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute nlb rule target remove --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --ip TARGET_IP --port TARGET_PORT
```

