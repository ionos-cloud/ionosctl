---
description: "Create or replace the bucket policy"
---

# ObjectStorageBucketPolicyPut

## Usage

```text
ionosctl object-storage bucket policy put [flags]
```

## Aliases

For `bucket` command:

```text
[b]
```

For `policy` command:

```text
[pol]
```

For `put` command:

```text
[p]
```

## Description

Create or replace the bucket's access policy. A bucket has at most one policy; this REPLACES it entirely.

Provide the policy as a JSON file via --json-properties. The document uses S3/IAM policy syntax: a "Version" (use "2012-10-17") and a "Statement" array. Each statement has:
  Effect     "Allow" or "Deny".
  Principal  Who it applies to. {"AWS": ["*"]} means everyone, including anonymous/public access.
  Action     S3 operations, e.g. "s3:GetObject", "s3:PutObject", "s3:ListBucket".
  Resource   ARNs of the bucket ("arn:aws:s3:::my-bucket") and/or its objects ("arn:aws:s3:::my-bucket/*"). Replace BUCKET_NAME in the example with your bucket.
  Condition  Optional extra constraints.

Granting Principal "*" with s3:GetObject makes objects publicly readable - but only if public-access-block does not block it (public-access-block wins). The "s3:"/"arn:aws:s3:::" tokens are the required S3-compatible format, not AWS-specific.

Run with --json-properties-example to print a ready-to-edit public-read template.

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Sid Effect Action Resource Principal Condition]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --json-properties string    Path to a JSON file with the bucket policy (IAM-style: Version + Statement[]). Replaces the existing policy
      --json-properties-example   Print an example bucket policy JSON and exit without contacting the API
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
# Apply a bucket policy from a file
ionosctl object-storage bucket policy put --name my-bucket --json-properties policy.json

# Print an example (public-read) policy to adapt
ionosctl object-storage bucket policy put --json-properties-example
```

