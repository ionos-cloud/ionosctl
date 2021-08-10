---
description: Update a User
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

Note: The password attribute is immutable. It is not allowed in update requests. It is recommended that the new User log into the DCD and change their password.

Required values to run command:

* User Id

## Options

```text
      --admin               Assigns the User to have administrative rights
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -e, --email string        The email for the User
      --first-name string   The first name for the User
  -f, --force               Force command to execute without user input
      --force-secure-auth   Indicates if secure (two-factor) authentication should be forced for the User
  -h, --help                help for update
      --last-name string    The last name for the User
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -i, --user-id string      The unique User Id (required)
  -v, --verbose             see step by step process when running a command
```

## Examples

```text
ionosctl user update --user-id USER_ID --admin=true
```

