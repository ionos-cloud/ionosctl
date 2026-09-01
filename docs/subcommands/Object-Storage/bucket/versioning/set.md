---
description: "Enable or suspend versioning on a bucket"
---

# ObjectStorageBucketVersioningSet

## Usage

```text
ionosctl object-storage bucket versioning set [flags]
```

## Aliases

For `bucket` command:

```text
[b]
```

For `versioning` command:

```text
[ver]
```

For `set` command:

```text
[s]
```

## Description

Set a bucket's versioning status to Enabled or Suspended.

Enabled starts keeping a distinct version for every overwrite/delete, protecting against accidental data loss. Suspended stops creating new versions from that point on but does NOT delete versions already stored, and does not return the bucket to a truly unversioned state. There is no way to fully disable versioning once it has been enabled.

Existing versions keep incurring storage cost; pair versioning with a lifecycle NoncurrentVersionExpiration rule if you want old versions cleaned up automatically. Versioning is also required for Object Lock.

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Name Versioning]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
  -n, --name string       Name of the bucket (required)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
      --status string     Target versioning status. 'Enabled' keeps a version per overwrite/delete; 'Suspended' stops new versions but retains existing ones (cannot be fully disabled once enabled) (required)
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Enable versioning
ionosctl object-storage bucket versioning set --name my-bucket --status Enabled

# Suspend versioning (existing versions are retained)
ionosctl object-storage bucket versioning set --name my-bucket --status Suspended
```

