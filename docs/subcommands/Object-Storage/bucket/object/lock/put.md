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

Set the DEFAULT retention that will be applied to every new object version uploaded to the bucket.

The bucket must already have Object Lock enabled - i.e. it was created with --object-lock (it cannot be turned on after the fact). This command marks ObjectLockEnabled=Enabled and installs a default retention rule.

Choose a mode and exactly one period:
  --mode   GOVERNANCE (users with bypass permission may shorten/remove) or COMPLIANCE (immutable for everyone, even root, until it expires - irreversible).
  --days   Retention length in days,  OR
  --years  Retention length in years. Provide exactly one of --days / --years; supplying both, or neither, is an error.

This default applies going forward to newly written versions; it does not retroactively lock objects that already exist, and per-object retention can still be set at upload time. Choose COMPLIANCE with care: once an object is locked in COMPLIANCE mode you cannot delete it (or the bucket) until every locked version's retention has elapsed.

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ObjectLockEnabled RetentionMode RetentionDays RetentionYears]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --days int32        Default retention period, in days, that new object versions stay locked. Provide exactly one of --days or --years
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
      --mode string       Default retention mode. GOVERNANCE: shortenable/removable by users with bypass permission. COMPLIANCE: immutable until expiry, even for the account root (required)
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
      --years int32       Default retention period, in years, that new object versions stay locked. Provide exactly one of --days or --years
```

## Examples

```text
# 30-day GOVERNANCE default (can be bypassed by privileged users)
ionosctl object-storage bucket object-lock put --name my-bucket --mode GOVERNANCE --days 30

# 1-year COMPLIANCE default (immutable, cannot be shortened or bypassed)
ionosctl object-storage bucket object-lock put --name my-bucket --mode COMPLIANCE --years 1
```

