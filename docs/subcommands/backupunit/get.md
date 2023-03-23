---
description: Get a BackupUnit
---

# BackupunitGet

## Usage

```text
ionosctl backupunit get [flags]
```

## Aliases

For `backupunit` command:

```text
[b backup]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific BackupUnit.

Required values to run command:

* BackupUnit Id

## Options

```text
  -i, --backupunit-id string   The unique BackupUnit Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --no-headers             When using text output, don't print headers
```

## Examples

```text
ionosctl backupunit get --backupunit-id BACKUPUNIT_ID
```

