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

Create a contract-owned bucket

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
  -l, --location string   Location of the resource to operate on. Can be one of: eu-central-3, eu-central-4, us-central-1 (default "eu-central-3")
  -n, --name string       Name of the bucket to create (required)
      --no-headers        Don't print table headers when table output is used
      --object-lock       Enable Object Lock on the new bucket (cannot be changed after creation)
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl object-storage bucket create --name my-bucket
ionosctl object-storage bucket create --name my-bucket --location eu-central-3
ionosctl object-storage bucket create --name my-bucket --object-lock
```

