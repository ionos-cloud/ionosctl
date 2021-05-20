---
description: Delete a BackupUnit
---

# BackupunitDelete

## Usage

```text
ionosctl backupunit delete [flags]
```

## Aliases

For `backupunit` command:
```text
[b]
```

## Description

Use this command to delete a BackupUnit. Deleting a BackupUnit is a dangerous operation. A successful DELETE will remove the backup plans inside a BackupUnit, ALL backups associated with the BackupUnit, the backup user and finally the BackupUnit itself.

Required values to run command:

* BackupUnit Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --backupunit-id string   The unique BackupUnit Id (required)
  -C, --cols strings           Set of columns to be printed on output 
                               Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                  Force command to execute without user input
  -h, --help                   help for delete
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for BackupUnit deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for BackupUnit deletion to be executed
```

## Examples

```text
ionosctl backupunit delete --backupunit-id 9fa48167-6375-4d93-b33c-e1ba3f461c17
Warning: Are you sure you want to delete backup unit (y/N) ? 
y
RequestId: fa00ba7e-426d-4460-9ec4-8b480bf5b17f
Status: Command backupunit delete has been successfully executed
```

