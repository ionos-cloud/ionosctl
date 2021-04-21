---
description: Delete a BackupUnit
---

# Delete

## Usage

```text
ionosctl backupunit delete [flags]
```

## Description

Use this command to delete a BackupUnit. Deleting a BackupUnit is a dangerous operation. A successful DELETE will remove the backup plans inside a BackupUnit, ALL backups associated with the BackupUnit, the backup user and finally the BackupUnit itself.

Required values to run command:

* BackupUnit Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --backupunit-id string   The unique BackupUnit Id [Required flag]
      --cols strings           Columns to be printed in the standard output (default [BackupUnitId,Name,Email])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                   help for delete
      --ignore-stdin           Force command to execute without user input
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --timeout int            Timeout option for BackupUnit to be deleted [seconds] (default 60)
      --wait                   Wait for BackupUnit to be deleted
```

## Examples

```text
ionosctl backupunit delete --backupunit-id 9fa48167-6375-4d93-b33c-e1ba3f461c17
Warning: Are you sure you want to delete backup unit (y/N) ? 
y
RequestId: fa00ba7e-426d-4460-9ec4-8b480bf5b17f
Status: Command backupunit delete has been successfully executed
```

