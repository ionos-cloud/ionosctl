---
description: Update a BackupUnit
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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
  -i, --backupunit-id string   The unique BackupUnit Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -e, --email string           The e-mail address you want to update for the BackupUnit
  -f, --force                  Force command to execute without user input
  -h, --help                   help for update
  -o, --output string          Desired output format [text|json] (default "text")
  -p, --password string        Alphanumeric password you want to update for the BackupUnit
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for BackupUnit update [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for BackupUnit update to be executed
```

## Examples

```text
ionosctl backupunit update --backupunit-id BACKUPUNIT_ID --email EMAIL
```

