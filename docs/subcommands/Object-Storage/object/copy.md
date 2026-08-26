---
description: "Server-side copy an object to a new key or bucket"
---

# ObjectStorageObjectCopy

## Usage

```text
ionosctl object-storage object copy [flags]
```

## Aliases

For `object-storage` command:

```text
[os]
```

For `object` command:

```text
[obj]
```

For `copy` command:

```text
[cp]
```

## Description

Copy an object to a new key, within the same bucket or into another bucket.

The copy happens entirely SERVER-SIDE: the bytes are never downloaded to your machine and re-uploaded. Both --name (the destination bucket) and --copy-source (the source) must live in the same region/endpoint.

--copy-source names the object being copied and must be in the form /source-bucket/source-key (a leading slash, then the bucket, then the key). The destination is --name plus --key. Copying onto an existing key overwrites it; on a versioning-enabled destination bucket the copy becomes a new version.

Common uses: rename or move an object (copy to the new key, then delete the old one), duplicate an object into another bucket, or restore an old version to the current key. The command prints the new object's ETag and LastModified on success.

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ETag LastModified]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --copy-source string   The source object to copy, in the form /source-bucket/source-key (leading slash required) (required)
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -k, --key string           Key to store the copy under in the destination bucket (overwrites if it already exists) (required)
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
  -n, --name string          Destination bucket to copy into (must be in the same region as the source) (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Copy within the same bucket (e.g. to "rename" a key - delete the old one afterwards)
ionosctl object-storage object copy --name my-bucket --key photos/renamed.jpg --copy-source /my-bucket/photos/image.jpg

# Copy across buckets
ionosctl object-storage object copy --name backup-bucket --key photos/image.jpg --copy-source /my-bucket/photos/image.jpg
```

