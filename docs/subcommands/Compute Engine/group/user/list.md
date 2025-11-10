---
description: "List Users from a Group"
---

# GroupUserList

## Usage

```text
ionosctl group user list [flags]
```

## Aliases

For `group` command:

```text
[g]
```

For `user` command:

```text
[u]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Users from a specific Group.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [firstname lastname email administrator forceSecAuth secAuthActive s3CanonicalUserId active]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId lastLogin]

## Options

```text
  -u, --api-url string    Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32       Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings   Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
      --group-id string   The unique Group Id (required)
  -h, --help              Print usage
      --limit int         Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Pagination offset: Number of items to skip before starting to collect the results
      --order-by string   Limits results to those containing a matching value for a specific property
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl group user list --group-id GROUP_ID
```

