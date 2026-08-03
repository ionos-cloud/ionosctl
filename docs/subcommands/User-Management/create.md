---
description: "Create a new User on the contract"
---

# UserCreate

## Usage

```text
ionosctl compute user create [flags]
```

## Aliases

For `user` command:

```text
[u]
```

For `create` command:

```text
[c]
```

## Description

Create a new User on your contract. First name, last name, email and password are required; the email address is the User's login and must be unique across IONOS Cloud.

A newly created User has NO permissions by default. Give them access either by making them an administrator (--administrator, full contract access) or - the usual approach - by adding them to one or more Groups afterwards with `ionosctl compute group user add`, so they inherit those groups' privileges.

Required values to run a command:

* First Name
* Last Name
* Email
* Password

## Options

```text
      --admin               Make the User a contract administrator with full access to everything on the contract, bypassing all group privileges. Leave false for a normal User whose access comes from group membership
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -e, --email string        The User's email address. This is the login identity and must be unique across IONOS Cloud (required)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
      --first-name string   The User's first name (required)
  -f, --force               Force command to execute without user input
      --force-secure-auth   Force two-factor (secure) authentication for this User: they must set up 2FA before they can sign in
  -h, --help                Print usage
      --last-name string    The User's last name (required)
      --limit int           Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -p, --password string     The User's initial password (at least 5 characters). The User can change it after first login (required)
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a standard user (grant permissions later via group membership)
ionosctl compute user create --first-name Jane --last-name Doe --email jane.doe@example.com --password 's3cr3tPw'

# Create a full contract administrator who also must use two-factor auth
ionosctl compute user create --first-name Admin --last-name User --email admin@example.com --password 's3cr3tPw' --administrator --force-secure-auth
```

