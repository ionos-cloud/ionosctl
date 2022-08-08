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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
  -i, --backupunit-id string   The unique BackupUnit Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int              Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for BackupUnit deletion [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait-for-request       Wait for the Request for BackupUnit deletion to be executed
```

## Examples

```text
ionosctl backupunit delete --backupunit-id BACKUPUNIT_ID
```

