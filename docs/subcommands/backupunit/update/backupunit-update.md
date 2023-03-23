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
  -i, --backupunit-id string   The unique BackupUnit Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -e, --email string           The e-mail address you want to update for the BackupUnit
  -p, --password string        Alphanumeric password you want to update for the BackupUnit
  -t, --timeout int            Timeout option for Request for BackupUnit update [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for BackupUnit update to be executed
```

## Examples

```text
ionosctl backupunit update --backupunit-id BACKUPUNIT_ID --email EMAIL
```

