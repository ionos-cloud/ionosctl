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
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl backupunit list
```

