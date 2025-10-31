---
description: "Update a User"
---

# UserUpdate

## Usage

```text
ionosctl user update [flags]
```

## Aliases

For `user` command:

```text
[u]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update details about a specific User including their privileges.

Required values to run command:

* User Id

## Options

```text
      --admin               Assigns the User to have administrative rights. E.g.: --admin=true, --admin=false
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -e, --email string        The email for the User
      --first-name string   The first name for the User
  -f, --force               Force command to execute without user input
      --force-secure-auth   Indicates if secure (two-factor) authentication should be forced for the User. E.g.: --force-secure-auth=true, --force-secure-auth=false
  -h, --help                Print usage
      --last-name string    The last name for the User
      --limit int           Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -p, --password string     The password for the User (must be at least 5 characters long)
  -q, --quiet               Quiet output
  -i, --user-id string      The unique User Id (required)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl user update --user-id USER_ID --admin=true
```

