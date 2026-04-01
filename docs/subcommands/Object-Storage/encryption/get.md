---
description: "Get the default encryption configuration for a bucket"
---

# ObjectStorageEncryptionGet

## Usage

```text
ionosctl object-storage encryption get [flags]
```

## Aliases

For `object-storage` command:

```text
[os]
```

For `encryption` command:

```text
[enc]
```

For `get` command:

```text
[g]
```

## Description

Get the default encryption configuration for a bucket

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [SSEAlgorithm]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -n, --name string       Name of the bucket (required)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl object-storage encryption get --name my-bucket
```

