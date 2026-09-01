---
description: "Turn an object's legal hold ON or OFF"
---

# ObjectStorageObjectLegalHoldPut

## Usage

```text
ionosctl object-storage object legal-hold put [flags]
```

## Aliases

For `object` command:

```text
[obj]
```

For `legal-hold` command:

```text
[lh]
```

For `put` command:

```text
[p]
```

## Description

Turn the legal hold on an object ON or OFF via --status.

While the hold is ON the object (version) cannot be deleted or overwritten, with no expiry date and no bypass - it stays locked until this same command sets it OFF. This is independent of any retention lock; the object remains protected while either is active.

Requires a bucket created with Object Lock enabled (it cannot be turned on for an existing bucket). On versioned buckets, pass --version-id to hold a specific version.

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Status]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -k, --key string          Key of the object to hold or release (required)
      --limit int           Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
  -n, --name string         Name of the Object-Lock-enabled bucket holding the object (required)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
      --status string       ON to place the legal hold (locks the object indefinitely), OFF to release it (required)
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
      --version-id string   Apply the hold to this specific object version instead of the current one (versioned buckets only)
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Place a legal hold on an object
ionosctl object-storage object legal-hold put --name my-bucket --key my-object --status ON

# Release the legal hold
ionosctl object-storage object legal-hold put --name my-bucket --key my-object --status OFF
```

