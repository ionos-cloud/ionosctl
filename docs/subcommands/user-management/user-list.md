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
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -F, --filters strings   Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -M, --max-results int   The maximum number of elements to return
      --no-headers        When using text output, don't print headers
      --order-by string   Limits results to those containing a matching value for a specific property
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl user list
```

