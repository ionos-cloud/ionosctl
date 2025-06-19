---
description: "Update a BackupUnit"
---

# BackupunitUpdate

## Usage

```text
ionosctl backupunit update [flags]
```

## Aliases

For `backupunit` command:

```text
[b backup]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update details about a specific BackupUnit. The password and the email may be updated.

Required values to run command:

* BackupUnit Id

## Options

```text
  -i, --backupunit-id string   The unique BackupUnit Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -e, --email string           The e-mail address you want to update for the BackupUnit
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -p, --password string        Alphanumeric password you want to update for the BackupUnit
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for BackupUnit update [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl backupunit update --backupunit-id BACKUPUNIT_ID --email EMAIL
```

