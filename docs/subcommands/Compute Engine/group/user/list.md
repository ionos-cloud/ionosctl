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
* filter by property: [active administrator email firstname forceSecAuth lastname s3CanonicalUserId secAuthActive]
* filter by metadata: [createdDate etag lastLogin]

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
      --group-id string     The unique Group Id (required)
  -h, --help                Print usage
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
      --order-by string     Limits results to those containing a matching value for a specific property
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Print step-by-step process when running command
```

## Examples

```text
ionosctl group user list --group-id GROUP_ID
```

