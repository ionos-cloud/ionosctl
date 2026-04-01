---
description: "List objects in a bucket"
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

List objects in a bucket

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Key ContentType ContentLength LastModified ETag]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
      --max-keys int32    Maximum number of objects to return (0 for no limit) (default 1000)
  -n, --name string       Name of the bucket (required)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -p, --prefix string     Filter objects by key prefix (e.g. photos/)
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl object-storage object list --name my-bucket
ionosctl object-storage object list --name my-bucket --prefix photos/ --max-keys 100
```

