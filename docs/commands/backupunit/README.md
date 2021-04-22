---
description: BackupUnit Operations
---

# BackupUnit

## Usage

```text
ionosctl backupunit [command]
```

## Description

The sub-command of `ionosctl backupunit` allows you to list, get, create, update, delete BackupUnits under your account.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [BackupUnitId,Name,Email])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for backupunit
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl backupunit create](create.md) | Create a BackupUnit |
| [ionosctl backupunit delete](delete.md) | Delete a BackupUnit |
| [ionosctl backupunit get](get.md) | Get a BackupUnit |
| [ionosctl backupunit get-sso-url](get-sso-url.md) | Get BackupUnit SSO URL |
| [ionosctl backupunit list](list.md) | List BackupUnits |
| [ionosctl backupunit update](update.md) | Update a BackupUnit |

