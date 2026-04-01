---
description: "Delete an object"
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

Delete an object

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Key ContentType ContentLength LastModified ETag]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -k, --key string          Object key to delete (required)
      --limit int           Maximum number of items to return per request (default 50)
  -n, --name string         Name of the bucket (required)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
      --version-id string   Version ID to delete a specific version
```

## Examples

```text
ionosctl object-storage object delete --name my-bucket --key photos/image.jpg
ionosctl object-storage object delete --name my-bucket --key photos/image.jpg --version-id abc123 -f
```

