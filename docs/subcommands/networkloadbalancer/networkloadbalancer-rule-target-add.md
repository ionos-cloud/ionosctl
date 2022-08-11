---
description: Add a Network Load Balancer Forwarding Rule Target
---

# NetworkloadbalancerRuleTargetAdd

## Usage

```text
ionosctl networkloadbalancer rule target add [flags]
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

Use this command to add a Forwarding Rule Target in a specified Network Load Balancer Forwarding Rule. You can also set Health Check Settings for Forwarding Rule Target. The Check parameter for Health Check Settings specifies whether the target VM's health is checked. If turned off, a target VM is always considered available. If turned on, the target VM is available when accepting periodic TCP connections, to ensure that it is really able to serve requests. The address and port to send the tests to are those of the target VM. The health check only consists of a connection attempt.

Regarding the Weight parameter, this parameter is used to adjust the target VM's weight relative to other target VMs. All target VMs will receive a load proportional to their weight relative to the sum of all weights, so the higher the weight, the higher the load. The default weight is 1, and the maximal value is 256. A value of 0 means the target VM will not participate in load-balancing but will still accept persistent connections. If this parameter is used to distribute the load according to target VM's capacity, it is recommended to start with values which can both grow and shrink, for instance between 10 and 100 to leave enough room above and below for later adjustments.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id
* Target Ip
* Target Port

## Options

```text
  -u, --api-url string                  Override default host url (default "https://api.ionos.com")
      --check                           [Health Check] Check specifies whether the target VM's health is checked (default true)
      --check-interval int              [Health Check] CheckInterval determines the duration (in milliseconds) between consecutive health checks (default 2000)
      --cols strings                    Set of columns to be printed on output 
                                        Available columns: [TargetIp TargetPort Weight Check CheckInterval Maintenance] (default [TargetIp,TargetPort,Weight,Check,CheckInterval,Maintenance])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string            The unique Data Center Id (required)
  -D, --depth int                       Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                           Force command to execute without user input
  -h, --help                            Print usage
      --ip string                       IP of a balanced target VM (required)
      --maintenance                     [Health Check]  Maintenance specifies if a target VM should be marked as down, even if it is not
      --networkloadbalancer-id string   The unique NetworkLoadBalancer Id (required)
  -o, --output string                   Desired output format [text|json] (default "text")
  -P, --port string                     Port of the balanced target service. Range: 1 to 65535 (required)
  -q, --quiet                           Quiet output
      --rule-id string                  The unique ForwardingRule Id (required)
  -t, --timeout int                     Timeout option for Request for Forwarding Rule Target creation [seconds] (default 300)
  -v, --verbose                         Print step-by-step process when running command
  -w, --wait-for-request                Wait for the Request for Forwarding Rule Target creation to be executed
  -W, --weight int                      Weight parameter is used to adjust the target VM's weight relative to other target VMs. Maximum: 256 (default 1)
```

## Examples

```text
ionosctl networkloadbalancer rule target add --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --target-ip TARGET_IP --target-port TARGET_PORT -w
```

