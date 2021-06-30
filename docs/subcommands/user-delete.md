---
description: Blacklists the User, disabling them
---

# UserDelete

## Usage

```text
ionosctl user delete [flags]
```

## Aliases

For `user` command:

```text
[u]
```

For `delete` command:

```text
[d]
```

## Description

This command blacklists the User, disabling them. The User is not completely purged, therefore if you anticipate needing to create a User with the same name in the future, we suggest renaming the User before you delete it.

Required values to run command:

* User Id

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for delete
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -i, --user-id string   The unique User Id (required)
```

## Examples

```text
ionosctl user delete --user-id USER_ID --force
```

