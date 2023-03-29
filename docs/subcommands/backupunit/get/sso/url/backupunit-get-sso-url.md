---
description: Get BackupUnit SSO URL
---

# BackupunitGetSsoUrl

## Usage

```text
ionosctl backupunit get-sso-url [flags]
```

## Aliases

For `backupunit` command:

```text
[b backup]
```

## Description

Use this command to access the GUI with a Single Sign On URL that can be retrieved from the Cloud API using this request. If you copy the entire value returned and paste it into a browser, you will be logged into the BackupUnit GUI.

Required values to run command:

* BackupUnit Id

## Options

```text
  -i, --backupunit-id string   The unique BackupUnit Id (required)
```

## Examples

```text
ionosctl backupunit get-sso-url --backupunit-id BACKUPUNIT_ID
```

