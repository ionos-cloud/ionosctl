---
description: List BackupUnits
---

# BackupunitList

## Usage

```text
ionosctl backupunit list [flags]
```

## Aliases

For `backupunit` command:
```text
[b backup]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to get a list of existing BackupUnits available on your account.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl backupunit list 
BackupUnitId                           Name          Email
9fa48167-6375-4d93-b33c-e1ba3f461c17   test1234567   testrandom20@ionos.com
```

