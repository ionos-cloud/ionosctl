---
description: "Set or extend the WORM retention lock on an object"
---

# ObjectStorageObjectRetentionPut

## Usage

```text
ionosctl object-storage object retention put [flags]
```

## Aliases

For `object` command:

```text
[obj]
```

For `retention` command:

```text
[ret]
```

For `put` command:

```text
[p]
```

## Description

Place an Object Lock retention lock on an object, protecting it from deletion and overwrite until --retain-until-date.

--mode selects how strict the lock is: GOVERNANCE can be overridden by users with the bypass permission, COMPLIANCE can never be overridden before the date (see the "retention" group help).

The retain-until date is always in the future. You can freely EXTEND an existing lock (set a later date) in either mode. SHORTENING or removing a lock is only possible for GOVERNANCE mode and requires --bypass-governance-retention plus the bypass permission; COMPLIANCE-mode dates can never be reduced.

Requires a bucket created with Object Lock enabled - it cannot be turned on for an existing bucket. On versioned buckets, pass --version-id to lock a specific version.

## Options

```text
  -u, --api-url string                Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'objectstorage' and env var 'IONOS_API_URL' (default "https://s3.%s.ionoscloud.com")
      --bypass-governance-retention   Required (with the bypass permission) to shorten or replace an existing GOVERNANCE-mode lock. No effect on COMPLIANCE mode
      --cols strings                  Set of columns to be printed on output 
                                      Available columns: [Mode RetainUntilDate]
  -c, --config string                 Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                     Level of detail for response objects (default 1)
  -F, --filters strings               Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                         Force command to execute without user input
  -h, --help                          Print usage
  -k, --key string                    Key of the object to lock (required)
      --limit int                     Maximum number of items to return per request (default 50)
  -l, --location string               Location of the resource to operate on. When unset, list commands query all locations. Can be one of: eu-central-3, eu-central-4, us-central-1. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint. Defaults to eu-central-3
      --mode string                   Retention strictness. GOVERNANCE: overridable by users with the bypass permission. COMPLIANCE: cannot be shortened or removed by anyone before the date (required)
  -n, --name string                   Name of the Object-Lock-enabled bucket holding the object (required)
      --no-headers                    Don't print table headers when table output is used
      --offset int                    Number of items to skip before starting to collect the results
      --order-by string               Property to order the results by
  -o, --output string                 Desired output format [text|json|api-json] (default "text")
      --query string                  JMESPath query string to filter the output
  -q, --quiet                         Quiet output
      --retain-until-date string      Future date until which the object stays locked, in RFC 3339 format (e.g. 2026-01-01T00:00:00Z). The lock lapses automatically once this passes (required)
  -t, --timeout int                   Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count                 Increase verbosity level [-v, -vv, -vvv]
      --version-id string             Apply the retention to this specific object version instead of the current one (versioned buckets only)
  -w, --wait                          Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Lock an object in GOVERNANCE mode until a date (RFC 3339)
ionosctl object-storage object retention put --name my-bucket --key my-object --mode GOVERNANCE --retain-until-date 2026-01-01T00:00:00Z

# Shorten/replace an existing GOVERNANCE lock (requires bypass permission)
ionosctl object-storage object retention put --name my-bucket --key my-object --mode GOVERNANCE --retain-until-date 2026-01-01T00:00:00Z --bypass-governance-retention
```

