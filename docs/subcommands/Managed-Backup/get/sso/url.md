---
description: "Get BackupUnit SSO URL"
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
  -u, --api-url string         Override default host URL. Preferred over the config file override 'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --backupunit-id string   The unique BackupUnit Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl backupunit get-sso-url --backupunit-id BACKUPUNIT_ID
```

