---
description: "Update a BackupUnit"
---

# BackupunitUpdate

## Usage

```text
ionosctl compute backupunit update [flags]
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

Use this command to update a BackupUnit. Only the --password and --email may be changed; the --name (backup login) is immutable and cannot be updated here (to rename, delete and recreate the unit).

Changing --password rotates the backup agent login secret. Like on create, the password is never returned by the API, so record any new value you set.

Required values to run command:

* BackupUnit Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --backupunit-id string   The unique BackupUnit Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [BackupUnitId Name Email State]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int              Level of detail for response objects (default 1)
  -e, --email string           New e-mail address for backup service reports
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -p, --password string        New login secret for the backup agent. Write-only: never returned by the API, so record it
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Change the report e-mail
ionosctl compute backupunit update --backupunit-id BACKUPUNIT_ID --email newops@example.com

# Rotate the login password
ionosctl compute backupunit update --backupunit-id BACKUPUNIT_ID --password 'N3wS3cret!'
```

