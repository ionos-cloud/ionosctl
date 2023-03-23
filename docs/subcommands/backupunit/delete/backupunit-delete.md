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
[b backup]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a BackupUnit. Deleting a BackupUnit is a dangerous operation. A successful DELETE will remove the backup plans inside a BackupUnit, ALL backups associated with the BackupUnit, the backup user and finally the BackupUnit itself.

Required values to run command:

* BackupUnit Id

## Options

```text
  -a, --all                    Delete all BackupUnits.
  -i, --backupunit-id string   The unique BackupUnit Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -t, --timeout int            Timeout option for Request for BackupUnit deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for BackupUnit deletion to be executed
```

## Examples

```text
ionosctl backupunit delete --backupunit-id BACKUPUNIT_ID
```

