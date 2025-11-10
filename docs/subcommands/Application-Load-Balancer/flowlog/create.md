---
description: "Create an Application Load Balancer FlowLog"
---

# ApplicationloadbalancerFlowlogCreate

## Usage

```text
ionosctl applicationloadbalancer flowlog create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create an Application Load Balancer FlowLog in a specified Application Load Balancer.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Bucket Name

## Options

```text
  -a, --action string                       Specifies the traffic action pattern. (default "ALL")
  -u, --api-url string                      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --applicationloadbalancer-id string   The unique ApplicationLoadBalancer Id (required)
      --cols strings                        Set of columns to be printed on output 
                                            Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string                       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string                The unique Data Center Id (required)
  -D, --depth int32                         Controls the detail depth of the response objects. Max depth is 10.
  -d, --direction string                    Specifies the traffic direction pattern. (default "INGRESS")
  -f, --force                               Force command to execute without user input
  -h, --help                                Print usage
      --limit int                           Pagination limit: Maximum number of items to return per request (default 50)
  -n, --name string                         The name of the Application Load Balancer FlowLog. (default "Unnamed ALB Flow Log")
      --no-headers                          Don't print table headers when table output is used
      --offset int                          Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string                       Desired output format [text|json|api-json] (default "text")
      --query string                        JMESPath query string to filter the output
  -q, --quiet                               Quiet output
  -b, --s3bucket string                     S3 bucket name of an existing IONOS Cloud S3 bucket. (required)
  -t, --timeout int                         Timeout option for Request for Application Load Balancer FlowLog creation [seconds] (default 300)
  -v, --verbose count                       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request                    Wait for the Request for Application Load Balancer FlowLog creation to be executed
```

## Examples

```text
ionosctl applicationloadbalancer flowlog create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --action ACTION --name NAME --direction DIRECTION --s3bucket BUCKET_NAME
```

