---
description: Update a BackupUnit
---

# Update

## Usage

```text
ionosctl backupunit update [flags]
```

## Description

Use this command to update details about a specific BackupUnit. The password and the email may be updated.

Required values to run command:

* BackupUnit Id

## Options

```text
  -u, --api-url string               Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --backupunit-email string      The e-mail address you want to update for the BackupUnit
      --backupunit-id string         The unique BackupUnit Id [Required flag]
      --backupunit-password string   Alphanumeric password you want to update for the BackupUnit
      --cols strings                 Columns to be printed in the standard output (default [BackupUnitId,Name,Email])
  -c, --config string                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                         help for update
      --ignore-stdin                 Force command to execute without user input
  -o, --output string                Desired output format [text|json] (default "text")
  -q, --quiet                        Quiet output
      --timeout int                  Timeout option for BackupUnit to be updated [seconds] (default 60)
      --wait                         Wait for BackupUnit to be updated
```

## Examples

```text
ionosctl backupunit update --backupunit-id 9fa48167-6375-4d93-b33c-e1ba3f461c17 --backupunit-email testrandom22@ionos.com
BackupUnitId                           Name          Email
9fa48167-6375-4d93-b33c-e1ba3f461c17   test1234567   testrandom22@ionos.com
RequestId: a91fbce0-bb98-4be1-9d7f-90d3f6da8ffe
Status: Command backupunit update has been successfully executed
```

