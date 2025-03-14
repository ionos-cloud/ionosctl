---
description: "Add a Target to a Target Group"
---

# TargetgroupTargetAdd

## Usage

```text
ionosctl targetgroup target add [flags]
```

## Aliases

For `targetgroup` command:

```text
[tg]
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

Use this command to add a Target to a Target Group. You will need to provide the IP, the port and the weight. Weight parameter is used to adjust the target VM's weight relative to other target VMs. All target VMs will receive a load proportional to their weight relative to the sum of all weights, so the higher the weight, the higher the load. The default weight is 1, and the maximal value is 256. A value of 0 means the target VM will not participate in load-balancing but will still accept persistent connections. If this parameter is used to distribute the load according to target VM's capacity, it is recommended to start with values which can both grow and shrink, for instance between 10 and 100 to leave enough room above and below for later adjustments.

Health Check can also be set. The `--check` option specifies whether the target VM's health is checked. If turned off, a target VM is always considered available. If turned on, the target VM is available when accepting periodic TCP connections, to ensure that it is really able to serve requests. The address and port to send the tests to are those of the target VM. The health check only consists of a connection attempt.

You can wait for the Request to be executed using `--wait-for-request` or `-w` option.

Required values to run command:

* Target Group Id
* Target Ip
* Target Port

## Options

```text
  -u, --api-url string          Override default host url (default "https://api.ionos.com")
      --cols strings            Set of columns to be printed on output 
                                Available columns: [TargetIp TargetPort Weight HealthCheckEnabled MaintenanceEnabled] (default [TargetIp,TargetPort,Weight,HealthCheckEnabled,MaintenanceEnabled])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                   Force command to execute without user input
      --health-check-enabled    Makes the target available only if it accepts periodic health check TCP connection attempts; when turned off, the target is considered always available. The health check only consists of a connection attempt to the address and port of the target. Default is True. (default true)
  -h, --help                    Print usage
      --ip ip                   The IP of the balanced target VM. (required)
  -m, --maintenance-enabled     Maintenance mode prevents the target from receiving balanced traffic.
      --no-headers              Don't print table headers when table output is used
  -o, --output string           Desired output format [text|json|api-json] (default "text")
  -P, --port int                The port of the balanced target service; valid range is 1 to 65535. (required) (default 8080)
  -q, --quiet                   Quiet output
  -i, --targetgroup-id string   The unique Target Group Id (required)
  -t, --timeout int             Timeout option for Request for Target Group Target addition [seconds] (default 60)
  -v, --verbose                 Print step-by-step process when running command
  -w, --wait                    Polls the request continuously until the operation is completed 
  -W, --weight int              Traffic is distributed in proportion to target weight, relative to the combined weight of all targets. A target with higher weight receives a greater share of traffic. Valid range is 0 to 256 and default is 1; targets with weight of 0 do not participate in load balancing but still accept persistent connections. It is best use values in the middle of the range to leave room for later adjustments. (default 1)
```

## Examples

```text
ionosctl targetgroup target add --targetgroup-id TARGET_GROUP_ID --ip TARGET_IP --port TARGET_PORT
```

