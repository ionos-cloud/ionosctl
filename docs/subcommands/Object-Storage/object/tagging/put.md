---
description: "Create or replace the entire tag set on an object"
---

# ObjectStorageObjectTaggingPut

## Usage

```text
ionosctl object-storage object tagging put [flags]
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

For `put` command:

```text
[p]
```

## Description

Create or replace the tag set on an object.

This REPLACES the object's whole tag set (it is not a merge): any existing tags not present in the supplied file are removed. To keep existing tags, include them in the file. An object may hold at most 10 tags.

The tag set is supplied as a JSON file via --json-properties, containing a "TagSet" array of {"Key", "Value"} objects. Run the command with --json-properties-example to print a ready-to-edit template.

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
      --json-properties string    Path to a JSON file with the full tag set to apply (a "TagSet" array of {"Key","Value"} pairs). See --json-properties-example
      --json-properties-example   Print an example tag-set JSON to stdout and exit, without contacting the API
  -k, --key string                Key of the object to tag (required)
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
  -n, --name string               Name of the bucket holding the object (required)
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
      --version-id string         Tag this specific object version instead of the current one (versioned buckets only)
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl object-storage object tagging put --name my-bucket --key my-object --json-properties tags.json
ionosctl object-storage object tagging put --json-properties-example
```

