---
description: "Blacklists the User, disabling them"
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
  -a, --all              Delete all the Users.
  -u, --api-url string   Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32      Controls the detail depth of the response objects. Max depth is 10.
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -i, --user-id string   The unique User Id (required)
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl user delete --user-id USER_ID --force
```

