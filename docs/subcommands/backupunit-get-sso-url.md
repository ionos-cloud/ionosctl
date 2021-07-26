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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
  -i, --backupunit-id string   The unique BackupUnit Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                  Force command to execute without user input
  -h, --help                   help for get-sso-url
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
```

## Examples

```text
ionosctl backupunit get-sso-url --backupunit-id BACKUPUNIT_ID
```

