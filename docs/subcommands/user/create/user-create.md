---
description: Create a User under a particular contract
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
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -e, --email string        The email for the User (required)
      --first-name string   The first name for the User (required)
      --force-secure-auth   Indicates if secure (two-factor) authentication should be forced for the User
      --last-name string    The last name for the User (required)
  -p, --password string     The password for the User (must be at least 5 characters long) (required)
```

## Examples

```text
ionosctl user create --first-name NAME --last-name NAME --email EMAIL --password PASSWORD
```

