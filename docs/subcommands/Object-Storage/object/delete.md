---
description: "Delete an object or all objects from a bucket"
---

# ObjectStorageObjectDelete

## Usage

```text
ionosctl object-storage object delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Delete a single object by key, or all objects (including versions and delete markers) from a bucket using --all.

## Options

```text
  -a, --all                           Delete all objects in the bucket
  -u, --api-url string                Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --bypass-governance-retention   Bypass Governance-mode Object Lock restrictions to delete the object
      --cols strings                  Set of columns to be printed on output 
                                      Available columns: [Key ContentType ContentLength LastModified ETag]
  -c, --config string                 Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                     Level of detail for response objects (default 1)
  -F, --filters strings               Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                         Force command to execute without user input
  -h, --help                          Print usage
  -k, --key string                    Object key to delete
      --limit int                     Maximum number of items to return per request (default 50)
  -l, --location string               Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1
  -n, --name string                   Name of the bucket (required)
      --no-headers                    Don't print table headers when table output is used
      --offset int                    Number of items to skip before starting to collect the results
      --order-by string               Property to order the results by
  -o, --output string                 Desired output format [text|json|api-json] (default "text")
      --query string                  JMESPath query string to filter the output
  -q, --quiet                         Quiet output
  -t, --timeout int                   Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                 Increase verbosity level [-v, -vv, -vvv]
      --version-id string             Version ID to delete a specific version
  -w, --wait                          Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl object-storage object delete --name my-bucket --key photos/image.jpg
ionosctl object-storage object delete --name my-bucket --key photos/image.jpg --version-id abc123 -f
ionosctl object-storage object delete --name my-bucket --all -f
ionosctl object-storage object delete --name my-bucket --all --bypass-governance-retention -f
```

