---
description: "Download an object's bytes to a local file"
---

# ObjectStorageObjectGet

## Usage

```text
ionosctl object-storage object get [flags]
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

For `get` command:

```text
[g]
```

## Description

Download the bytes of an object to a local file.

By default the file is written to the current directory under the basename of the key (the part after the last "/", so photos/image.jpg is saved as image.jpg). Use --destination to write elsewhere or under a different name.

On a versioning-enabled bucket, the current (latest) version is downloaded unless you pass --version-id to fetch a specific historical version. Version IDs are shown by "object list" (with the appropriate columns) and in the S3 version listing.

To fetch only an object's metadata (size, content-type, ETag, last-modified) without downloading the bytes, use "object head" instead.

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Key ContentType ContentLength LastModified ETag]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -d, --destination string   Local file path to write the download to. Defaults to the basename of the key (the part after the last "/") in the current directory
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -k, --key string           Key (full name) of the object to download (required)
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
  -n, --name string          Name of the bucket holding the object (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
      --version-id string    Download this specific object version instead of the current one (versioned buckets only). Defaults to the latest version
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Download to ./image.jpg (basename of the key) in the current directory
ionosctl object-storage object get --name my-bucket --key photos/image.jpg

# Download to an explicit local path
ionosctl object-storage object get --name my-bucket --key photos/image.jpg --destination ./local-image.jpg

# Download a specific historical version
ionosctl object-storage object get --name my-bucket --key photos/image.jpg --version-id <version-id>
```

