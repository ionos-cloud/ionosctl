---
description: "Update an Application Load Balancer FlowLog"
---

# ApplicationloadbalancerFlowlogUpdate

## Usage

```text
ionosctl applicationloadbalancer flowlog update [flags]
```

## Aliases

For `applicationloadbalancer` command:

```text
[alb]
```

For `flowlog` command:

```text
[f fl]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified Application Load Balancer FlowLog from an Application Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Application Load Balancer FlowLog Id

## Options

```text
  -a, --action string                       Specifies the traffic action pattern.
  -u, --api-url string                      Override default host url (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -d, --direction string                    Specifies the traffic direction pattern.
  -i, --flowlog-id string                   The unique FlowLog Id (required)
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
  -n, --name string                         The name of the Application Load Balancer FlowLog.
      --no-headers                          Don't print table headers when table output is used
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
  -q, --quiet                               Quiet output
  -b, --s3bucket string                     S3 bucket name of an existing IONOS Cloud S3 bucket.
  -t, --timeout int                         Timeout option for Request for Application Load Balancer FlowLog update [seconds] (default 300)
  -v, --verbose                             Print step-by-step process when running command
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer FlowLog update to be executed
```

## Examples

```text
ionosctl applicationloadbalancer flowlog update --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FLOWLOG_ID --name NAME
```

