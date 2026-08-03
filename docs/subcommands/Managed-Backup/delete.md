---
description: "Delete a BackupUnit"
---

# BackupunitDelete

## Usage

```text
ionosctl compute backupunit delete [flags]
```

## Aliases

For `backupunit` command:

```text
[b backup]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a BackupUnit. This is a DESTRUCTIVE and irreversible operation: a successful delete removes the backup plans inside the unit, ALL backups stored in it, the backup login user, and finally the BackupUnit itself.

Because the name (backup login) is immutable, deleting is also the only way to "rename" a unit: delete and recreate under a new name (note the recreated unit starts empty).

Required values to run command:

* BackupUnit Id

## Options

```text
  -a, --all                    Delete all BackupUnits under the contract (each with its backups). Use instead of --backupunit-id
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --backupunit-id string   The unique BackupUnit Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [BackupUnitId Name Email State]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Delete one BackupUnit
ionosctl compute backupunit delete --backupunit-id BACKUPUNIT_ID

# Delete every BackupUnit under the contract
ionosctl compute backupunit delete --all
```

