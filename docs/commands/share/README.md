---
description: Resource Share Operations
---

# Share

## Usage

```text
ionosctl share [command]
```

## Description

The sub-commands of `ionosctl share` allow you to list, get, create, update, delete Resource Shares.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [ShareId,EditPrivilege,SharePrivilege,Type])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force            Force command to execute without user input
  -h, --help             help for share
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl share create](create.md) | Create a Resource Share for a Group |
| [ionosctl share delete](delete.md) | Delete a Resource Share from a Group |
| [ionosctl share get](get.md) | Get a Resource Share from a Group |
| [ionosctl share list](list.md) | List Resources Shares through a Group |
| [ionosctl share update](update.md) | Update a Resource Share from a Group |

