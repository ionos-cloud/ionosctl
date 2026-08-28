---
description: "Create a contract-owned bucket"
---

# ObjectStorageBucketCreate

## Usage

```text
ionosctl object-storage bucket create [flags]
```

## Aliases

For `object-storage` command:

```text
[os]
```

For `bucket` command:

```text
[b]
```

For `create` command:

```text
[c]
```

## Description

Create a new S3-compatible bucket owned by your contract.

The bucket name must be globally unique across the whole Object Storage service (not just your account) and follow S3 bucket-naming rules (3-63 chars, lowercase letters, numbers, dots and hyphens; DNS-compatible). The bucket is created in a single location and stays there for its lifetime.

--location selects the region the bucket is created in and sets its location constraint; it also decides which regional endpoint the request is signed for, so the two always stay in sync. When --location is omitted, the first Object Storage location (`eu-central-3`) is used.

--object-lock enables WORM (Write-Once-Read-Many) Object Lock on the bucket. This is a create-time-only decision: it CANNOT be turned on (or off) after the bucket exists. Enabling it also implicitly enables versioning. After creation, define the default retention with `bucket object-lock put`.

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Name CreationDate Region]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
  -n, --name string       Globally-unique bucket name (3-63 chars, lowercase, DNS-compatible). Must not already exist anywhere in the service (required)
      --no-headers        Don't print table headers when table output is used
      --object-lock       Enable WORM Object Lock at creation. Irreversible and implicitly enables versioning; set the default retention afterwards with 'bucket object-lock put'. Default: false
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a bucket in the default location
ionosctl object-storage bucket create --name my-bucket

# Create a bucket in a specific region
ionosctl object-storage bucket create --name my-bucket --location eu-central-3

# Create a bucket with WORM Object Lock enabled (irreversible; implies versioning)
ionosctl object-storage bucket create --name my-bucket --object-lock
```

