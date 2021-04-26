---
description: Private Cross-Connect Operations
---

# PrivateCrossConnect

## Usage

```text
ionosctl pcc [command]
```

## Description

The sub-command of `ionosctl pcc` allows you to list, get, create, update, delete Private Cross-Connect. To add Private Cross-Connect to a Lan, check the `ionosctl lan update` command.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force            Force command to execute without user input
  -h, --help             help for pcc
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl pcc create](create.md) | Create a Private Cross-Connect |
| [ionosctl pcc delete](delete.md) | Delete a Private Cross-Connect |
| [ionosctl pcc get](get.md) | Get a Private Cross-Connect |
| [ionosctl pcc get-peers](get-peers.md) | Get a Private Cross-Connect Peers |
| [ionosctl pcc list](list.md) | List Private Cross-Connects |
| [ionosctl pcc update](update.md) | Update a Private Cross-Connect |

