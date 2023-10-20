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
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -e, --email string        The email for the User
      --first-name string   The first name for the User
  -f, --force               Force command to execute without user input
      --force-secure-auth   Indicates if secure (two-factor) authentication should be forced for the User. E.g.: --force-secure-auth=true, --force-secure-auth=false
  -h, --help                Print usage
      --last-name string    The last name for the User
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -p, --password string     The password for the User (must be at least 5 characters long)
  -q, --quiet               Quiet output
  -i, --user-id string      The unique User Id (required)
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl user update --user-id USER_ID --admin=true
```

