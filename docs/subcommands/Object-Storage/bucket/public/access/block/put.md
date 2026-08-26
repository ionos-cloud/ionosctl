---
description: "Create or replace the public access block configuration for a bucket"
---

# ObjectStorageBucketPublicAccessBlockPut

## Usage

```text
ionosctl object-storage bucket public-access-block put [flags]
```

## Aliases

For `bucket` command:

```text
[b]
```

For `public-access-block` command:

```text
[pab]
```

For `put` command:

```text
[p]
```

## Description

Create or replace the bucket's Public Access Block. These guardrails override ACLs and bucket policies, so turning them on is the dependable way to force a bucket private regardless of what its policy/ACLs say.

Provide the configuration as a JSON file via --json-properties - a flat object with four booleans:
  BlockPublicAcls        Reject requests that would apply a public ACL.
  IgnorePublicAcls       Ignore public ACLs already set.
  BlockPublicPolicy      Reject bucket policies that grant public access.
  RestrictPublicBuckets  Limit access via an existing public policy to authorized principals.

Set all four to true to lock a bucket down completely. This is a full REPLACE of the configuration. Run with --json-properties-example to print a template with every guardrail enabled.

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [BlockPublicAcls IgnorePublicAcls BlockPublicPolicy RestrictPublicBuckets]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --json-properties string    Path to a JSON file with the four block flags (BlockPublicAcls, IgnorePublicAcls, BlockPublicPolicy, RestrictPublicBuckets). Replaces the existing configuration
      --json-properties-example   Print an example public access block configuration JSON and exit without contacting the API
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
# Fully block public access (all four guardrails on)
ionosctl object-storage bucket public-access-block put --name my-bucket --json-properties config.json

# Print an example configuration
ionosctl object-storage bucket public-access-block put --json-properties-example
```

