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

Required values to run command:

* User Id

## Options

```text
      --admin               Assigns the User to have administrative rights. E.g.: --admin=true, --admin=false
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -e, --email string        The email for the User
      --first-name string   The first name for the User
      --force-secure-auth   Indicates if secure (two-factor) authentication should be forced for the User. E.g.: --force-secure-auth=true, --force-secure-auth=false
      --last-name string    The last name for the User
  -p, --password string     The password for the User (must be at least 5 characters long)
  -i, --user-id string      The unique User Id (required)
```

## Examples

```text
ionosctl user update --user-id USER_ID --admin=true
```

