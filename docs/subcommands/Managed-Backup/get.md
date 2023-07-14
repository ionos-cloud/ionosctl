---
description: "Get a BackupUnit"
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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
  -i, --backupunit-id string   The unique BackupUnit Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             When using text output, don't print headers
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl backupunit get --backupunit-id BACKUPUNIT_ID
```

