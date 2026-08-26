---
description: "Upload a local file as an object"
---

# ObjectStorageObjectPut

## Usage

```text
ionosctl object-storage object put [flags]
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

For `put` command:

```text
[p]
```

## Description

Upload a local file into a bucket, stored under the key you choose with --key.

The key is the object's full name in the bucket, including any "/" separators you want to appear as a pseudo-folder path (e.g. photos/2025/image.jpg). The key does NOT have to match the source filename. If the key already exists it is overwritten; on a versioning-enabled bucket the overwrite creates a new version rather than destroying the old bytes.

--content-type sets the object's MIME type, which the service stores and returns as the Content-Type header on later downloads (it drives how browsers and clients interpret the bytes). If omitted, it is auto-detected from the source file extension, falling back to application/octet-stream when the extension is unknown.

The whole file is read from --source and uploaded as the object body.

## Options

```text
  -u, --api-url string        Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [Key ContentType ContentLength LastModified ETag]
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --content-type string   MIME type stored with the object and returned as its Content-Type on download (e.g. image/jpeg). Auto-detected from the --source file extension when omitted, defaulting to application/octet-stream if unknown
  -D, --depth int             Level of detail for response objects (default 1)
  -F, --filters strings       Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
  -k, --key string            Key to store the object under, i.e. its full name in the bucket. Use "/" to build a pseudo-folder path (e.g. photos/image.jpg). An existing key is overwritten (creating a new version on versioned buckets) (required)
      --limit int             Maximum number of items to return per request (default 50)
  -l, --location string       Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
  -n, --name string           Name of the destination bucket to upload into (required)
      --no-headers            Don't print table headers when table output is used
      --offset int            Number of items to skip before starting to collect the results
      --order-by string       Property to order the results by
  -o, --output string         Desired output format [text|json|api-json] (default "text")
      --query string          JMESPath query string to filter the output
  -q, --quiet                 Quiet output
  -s, --source string         Path to the local file whose bytes are uploaded as the object body (required)
  -t, --timeout int           Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count         Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                  Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Upload a file, letting the content-type be auto-detected from the extension
ionosctl object-storage object put --name my-bucket --key photos/image.jpg --source ./image.jpg

# Upload under a different key and set an explicit content-type
ionosctl object-storage object put --name my-bucket --key exports/report.json --source ./out.dat --content-type application/json
```

