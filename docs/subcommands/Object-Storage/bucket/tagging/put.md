---
description: "Create or replace the tagging configuration for a bucket"
---

# ObjectStorageBucketTaggingPut

## Usage

```text
ionosctl object-storage bucket tagging put [flags]
```

## Aliases

For `bucket` command:

```text
[b]
```

For `tagging` command:

```text
[tag]
```

For `put` command:

```text
[p]
```

## Description

Create or replace the bucket's tag set. This REPLACES all existing tags with the ones in the file - it is not a merge - so include every tag you want the bucket to keep.

Provide the tags as a JSON file via --json-properties: a top-level object with a "TagSet" array of {"Key": ..., "Value": ...} pairs. Tags are commonly used for cost allocation and grouping resources by environment, team or project.

Run with --json-properties-example to print a ready-to-edit template.

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Key Value]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --json-properties string    Path to a JSON file with the tag set ({"TagSet":[{"Key":...,"Value":...}]}). Replaces all existing tags
      --json-properties-example   Print an example tagging configuration JSON and exit without contacting the API
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
  -n, --name string               Name of the bucket (required)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Apply a set of tags from a file
ionosctl object-storage bucket tagging put --name my-bucket --json-properties tags.json

# Print an example tag set
ionosctl object-storage bucket tagging put --json-properties-example
```

