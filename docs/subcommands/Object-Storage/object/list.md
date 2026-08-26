---
description: "List objects (keys) in a bucket"
---

# ObjectStorageObjectList

## Usage

```text
ionosctl object-storage object list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

List the objects in a bucket, showing each key with its size, storage class, last-modified time and ETag.

Only CURRENT objects are listed - on a versioning-enabled bucket, older versions and delete markers are not shown here. Keys are returned in lexicographic (alphabetical) order.

--prefix restricts the listing to keys that start with a given string. Because "/" in a key is just a naming convention, a trailing-slash prefix like photos/ acts as a pseudo-folder filter.

--max-keys caps how many objects are returned; the command paginates transparently under the hood (in pages of up to 1000) until it has that many or the bucket is exhausted. Pass 0 to return every object in the bucket with no cap.

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Key Size LastModified StorageClass ETag]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
      --max-keys int32    Maximum number of objects to return; the command paginates transparently to reach it. Use 0 for no limit (list the entire bucket). Default: 1000 (default 1000)
  -n, --name string       Name of the bucket to list objects from (required)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -p, --prefix string     Only list keys beginning with this string. Use a trailing "/" (e.g. photos/) to browse a pseudo-folder
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# List up to 1000 objects (the default cap)
ionosctl object-storage object list --name my-bucket

# List only keys under the photos/ pseudo-folder, capped at 100
ionosctl object-storage object list --name my-bucket --prefix photos/ --max-keys 100

# List every object in the bucket (no cap)
ionosctl object-storage object list --name my-bucket --max-keys 0
```

