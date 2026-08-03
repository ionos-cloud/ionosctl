---
description: "Add a Network Load Balancer Forwarding Rule Target"
---

# NetworkloadbalancerRuleTargetAdd

## Usage

```text
ionosctl compute networkloadbalancer rule target add [flags]
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

For `add` command:

```text
[a]
```

## Description

Use this command to add a backend target to a forwarding rule. A target is a VM identified by --ip and --port on the NLB's target LAN; once added, the rule starts distributing connections to it according to the rule's balancing algorithm.

Weight (--weight): traffic is distributed in proportion to a target's weight relative to the sum of all targets' weights, so a higher weight means a higher share of connections. Default is 1, maximum is 256. A weight of 0 excludes the target from balancing but still lets it accept persistent connections. When sizing by capacity, start with mid-range values (e.g. 10-100) so you can adjust up or down later.

Health check (--check): when on (the default), the NLB periodically opens a TCP connection to the target's IP+port; the target is only considered available while it accepts these probes. When off, the target is always considered available. Use --check-interval to set the probe frequency, and --maintenance to force a target "down" (drain it) without removing it.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id
* Target Ip
* Target Port

## Options

```text
  -u, --api-url string                  Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --check                           [Health Check] When true, the target is only used while it accepts periodic TCP health probes; when false it is always considered available (default true)
      --check-interval int              [Health Check] Interval in milliseconds between consecutive TCP health probes to the target (default 2000)
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [TargetIp TargetPort Weight Check CheckInterval Maintenance]
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int                       Level of detail for response objects (default 1)
  -F, --filters strings                 Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --ip ip                           IP of the backend target VM on the NLB's target LAN (required)
      --limit int                       Maximum number of items to return per request (default 50)
      --maintenance                     [Health Check] When true, drains the target: it is treated as down and receives no balanced traffic, even if healthy
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
      --no-headers                      Don't print table headers when table output is used
      --offset int                      Number of items to skip before starting to collect the results
      --order-by string                 Property to order the results by
  -o, --output string                   Desired output format [text|json|api-json] (default "text")
  -P, --port string                     Port of the backend target service. Range: 1 to 65535 (required)
      --query string                    JMESPath query string to filter the output
  -q, --quiet                           Quiet output
      --rule-id string                  The unique ForwardingRule Id (required)
  -t, --timeout int                     Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                   Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                            Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
  -W, --weight int                      Share of traffic this target receives relative to the other targets' weights. Range: 0 to 256; 0 excludes it from balancing but still accepts persistent connections (default 1)
```

## Examples

```text
# Add a backend VM with default weight and health checks
ionosctl compute networkloadbalancer rule target add --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --ip 10.0.0.11 --port 80

# Add a higher-capacity backend (double share) with a faster health-check probe
ionosctl compute networkloadbalancer rule target add --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --ip 10.0.0.12 --port 80 --weight 200 --check-interval 1000
```

