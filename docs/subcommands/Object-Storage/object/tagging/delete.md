---
description: "Delete the tagging configuration for an object"
---

# ObjectStorageObjectTaggingDelete

## Usage

```text
ionosctl object-storage object tagging delete [flags]
```

## Aliases

For `object` command:

```text
[obj]
```

For `tagging` command:

```text
[tag]
```

For `delete` command:

```text
[d]
```

## Description

Delete the tagging configuration for an object

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Key Value]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -k, --key string          Object key (required)
      --limit int           Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. Defaults to eu-central-3
  -n, --name string         Name of the bucket (required)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
      --version-id string   Version ID of the object
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl object-storage object tagging delete --name my-bucket --key my-object
ionosctl object-storage object tagging delete --name my-bucket --key my-object -f
```

