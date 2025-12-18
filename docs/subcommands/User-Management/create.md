---
description: "Create a User under a particular contract"
---

# UserCreate

## Usage

```text
ionosctl user create [flags]
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

Use this command to create a User under a particular contract. You need to specify the firstname, lastname, email and password for the new User.

Required values to run a command:

* First Name
* Last Name
* Email
* Password

## Options

```text
      --admin               Assigns the User to have administrative rights
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -e, --email string        The email for the User (required)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
      --first-name string   The first name for the User (required)
  -f, --force               Force command to execute without user input
      --force-secure-auth   Indicates if secure (two-factor) authentication should be forced for the User
  -h, --help                Print usage
      --last-name string    The last name for the User (required)
      --limit int           Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -p, --password string     The password for the User (must be at least 5 characters long) (required)
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl user create --first-name NAME --last-name NAME --email EMAIL --password PASSWORD
```

