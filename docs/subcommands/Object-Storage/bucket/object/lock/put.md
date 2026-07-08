---
description: "Apply an Object Lock configuration to a bucket"
---

# ObjectStorageBucketObjectLockPut

## Usage

```text
ionosctl object-storage bucket object-lock put [flags]
```

## Aliases

For `bucket` command:

```text
[b]
```

For `object-lock` command:

```text
[ol]
```

For `put` command:

```text
[p]
```

## Description

Apply an Object Lock configuration to a bucket. The bucket must have been created with --object-lock enabled. Specify a default retention mode (GOVERNANCE or COMPLIANCE) and period (--days or --years, but not both).

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ObjectLockEnabled RetentionMode RetentionDays RetentionYears]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --days int32        Default retention period in days (mutually exclusive with --years)
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. Defaults to eu-central-3
      --mode string       Default retention mode: GOVERNANCE or COMPLIANCE (required)
  -n, --name string       Name of the bucket (required)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
      --years int32       Default retention period in years (mutually exclusive with --days)
```

## Examples

```text
ionosctl object-storage bucket object-lock put --name my-bucket --mode GOVERNANCE --days 30
ionosctl object-storage bucket object-lock put --name my-bucket --mode COMPLIANCE --years 1
```

