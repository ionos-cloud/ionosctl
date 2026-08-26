---
description: "Create or replace the lifecycle configuration for a bucket"
---

# ObjectStorageBucketLifecyclePut

## Usage

```text
ionosctl object-storage bucket lifecycle put [flags]
```

## Aliases

For `bucket` command:

```text
[b]
```

For `lifecycle` command:

```text
[lc]
```

For `put` command:

```text
[p]
```

## Description

Create or replace a bucket's lifecycle configuration. This fully REPLACES any existing rules with the "Rules" array in the file.

Provide the configuration as a JSON file via --json-properties. Each rule needs an "ID", a "Prefix" (which object keys it applies to; "" = all objects), a "Status" ("Enabled" or "Disabled"), and at least one action:
  Expiration.Days / Expiration.Date            Delete current objects after N days, or on a specific date.
  NoncurrentVersionExpiration.NoncurrentDays   Delete old versions N days after they become noncurrent (versioned buckets).
  AbortIncompleteMultipartUpload.DaysAfterInitiation   Abort unfinished multipart uploads after N days.

Note: on a versioned bucket, Expiration alone only adds a delete marker and does not reclaim storage - pair it with NoncurrentVersionExpiration. ionosctl computes and sends the required Content-MD5 for the request automatically.

Run with --json-properties-example to print a ready-to-edit template (expire objects under logs/ after 90 days, plus abort stale multipart uploads after 7 days).

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ID Prefix Status ExpirationDays ExpirationDate ExpiredObjectDeleteMarker NoncurrentDays AbortDays]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --json-properties string    Path to a JSON file with the lifecycle rules ({"Rules":[...]}). Replaces all existing rules
      --json-properties-example   Print an example lifecycle configuration JSON and exit without contacting the API
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
  -n, --name string               Name of the bucket (required)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Apply lifecycle rules from a file
ionosctl object-storage bucket lifecycle put --name my-bucket --json-properties lifecycle.json

# Print an example configuration
ionosctl object-storage bucket lifecycle put --json-properties-example
```

