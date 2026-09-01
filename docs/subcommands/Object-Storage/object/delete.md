---
description: "Delete a single object, or empty a whole bucket"
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

Delete one object by key, or empty the entire bucket with --all.

Versioning changes what "delete" means:
  - On a versioning-enabled bucket, deleting a key WITHOUT --version-id does not remove data: it inserts a "delete marker" that hides the key, while all prior versions remain recoverable.
  - Passing --version-id permanently removes that one specific version (this cannot be undone). Deleting a delete-marker version un-hides the key.
  - On a bucket without versioning, the object is simply removed.

--all empties the bucket: it deletes every current object AND every historical version AND every delete marker, so the bucket ends up truly empty. This is destructive and irreversible - it is guarded by a confirmation prompt (use -f/--force to skip it in scripts).

Object Lock: objects protected by a retention or legal hold cannot be deleted. GOVERNANCE-mode retention can be overridden with --bypass-governance-retention (requires the appropriate permission); COMPLIANCE-mode retention and active legal holds can never be bypassed.

## Options

```text
  -a, --all                           Empty the entire bucket: delete every object, version and delete marker. Destructive and irreversible
  -u, --api-url string                Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --bypass-governance-retention   Override GOVERNANCE-mode Object Lock retention so locked objects can be deleted (needs bypass permission). Has no effect on COMPLIANCE-mode retention or legal holds
      --cols strings                  Set of columns to be printed on output 
                                      Available columns: [Key ContentType ContentLength LastModified ETag]
  -c, --config string                 Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                     Level of detail for response objects (default 1)
  -F, --filters strings               Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                         Force command to execute without user input
  -h, --help                          Print usage
  -k, --key string                    Key of the object to delete. Mutually exclusive with --all
      --limit int                     Maximum number of items to return per request (default 50)
  -l, --location string               Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
  -n, --name string                   Name of the bucket to delete from (required)
      --no-headers                    Don't print table headers when table output is used
      --offset int                    Number of items to skip before starting to collect the results
      --order-by string               Property to order the results by
  -o, --output string                 Desired output format [text|json|api-json] (default "text")
      --query string                  JMESPath query string to filter the output
  -q, --quiet                         Quiet output
  -t, --timeout int                   Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                 Increase verbosity level [-v, -vv, -vvv]
      --version-id string             Permanently delete this specific version (irreversible). Without it, deleting on a versioned bucket only inserts a delete marker
  -w, --wait                          Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Delete one object (inserts a delete marker on a versioned bucket)
ionosctl object-storage object delete --name my-bucket --key photos/image.jpg

# Permanently delete one specific version
ionosctl object-storage object delete --name my-bucket --key photos/image.jpg --version-id <version-id> -f

# Empty the whole bucket (all objects, versions and delete markers)
ionosctl object-storage object delete --name my-bucket --all -f

# Empty a bucket, overriding GOVERNANCE-mode Object Lock protection
ionosctl object-storage object delete --name my-bucket --all --bypass-governance-retention -f
```

