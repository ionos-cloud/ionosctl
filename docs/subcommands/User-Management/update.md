---
description: "Update a User's details or admin/2FA status"
---

# UserUpdate

## Usage

```text
ionosctl compute user update [flags]
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

Update a User's profile fields (name, email, password) and their contract-level flags (administrator, forced two-factor auth). Only the fields you pass are changed.

Note: this command does NOT change which Groups a User belongs to - manage membership (and therefore inherited group privileges) with `ionosctl compute group user add/remove`. Toggling --administrator here is the one way to grant or revoke blanket contract-wide access directly on the User.

Required values to run command:

* User Id

## Options

```text
      --admin               Grant (true) or revoke (false) contract-administrator rights - full access to the whole contract, bypassing group privileges. E.g.: --admin=true, --admin=false
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -e, --email string        The User's email address (login identity). Must remain unique across IONOS CLOUD
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
      --first-name string   The User's first name
  -f, --force               Force command to execute without user input
      --force-secure-auth   Force (true) or stop forcing (false) two-factor authentication for this User. E.g.: --force-secure-auth=true, --force-secure-auth=false
  -h, --help                Print usage
      --last-name string    The User's last name
      --limit int           Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -p, --password string     Set a new password for the User (at least 5 characters)
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -i, --user-id string      The unique User Id (required)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Promote a user to contract administrator
ionosctl compute user update --user-id USER_ID --admin=true

# Demote an admin back to a normal user and require 2FA
ionosctl compute user update --user-id USER_ID --admin=false --force-secure-auth=true

# Reset a user's password
ionosctl compute user update --user-id USER_ID --password 'newS3cr3t'
```

