---
description: User Operations
---

# User

## Usage

```text
ionosctl user [command]
```

## Description

The sub-command of `ionosctl user` allows you to list, get, create, update, delete Users under your account. To add Users to a Group, check the `ionosctl group` commands.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for user
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl user create](create.md) | Create a User under a particular contract |
| [ionosctl user delete](delete.md) | Blacklists the User, disabling them |
| [ionosctl user get](get.md) | Get a User |
| [ionosctl user list](list.md) | List Users |
| [ionosctl user update](update.md) | Update a User |

