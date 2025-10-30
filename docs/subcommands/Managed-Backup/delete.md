---
description: "Delete a BackupUnit"
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
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --backupunit-id string   The unique BackupUnit Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for BackupUnit deletion [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl backupunit delete --backupunit-id BACKUPUNIT_ID
```

