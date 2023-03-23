---
description: List Users
---

# UserList

## Usage

```text
ionosctl user list [flags]
```

## Aliases

For `user` command:

```text
[u]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of existing Users available on your account.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [firstname lastname email administrator forceSecAuth secAuthActive s3CanonicalUserId active]
* filter by metadata: [etag createdDate lastLogin]

## Options

```text
      --cols strings        Set of columns to be printed on output 
                            Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          When using text output, don't print headers
      --order-by string     Limits results to those containing a matching value for a specific property
```

## Examples

```text
ionosctl user list
```

