---
description: "Create a BackupUnit"
---

# BackupunitCreate

## Usage

```text
ionosctl compute backupunit create [flags]
```

## Aliases

For `backupunit` command:

```text
[b backup]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a BackupUnit under your contract. A BackupUnit is the named storage container + login that the IONOS Managed Backup agent uses to store server backups.

You must supply --name, --email and --password.

Notes:

* --name becomes the backup login: it is concatenated with your contract number as CONTRACT_NUMBER-NAME, so it must be GLOBALLY UNIQUE across all IONOS contracts. It CANNOT be changed after creation (to rename, delete and recreate).
* --password is the login secret used to register the backup agent. It is WRITE-ONLY: the IONOS CLOUD API never returns it, so record it now. It can be changed later with `backupunit update`.
* --email receives service reports from the backup system and does NOT need to match your IONOS CLOUD API username. It can be changed later.
* After creation, log in to the backup console at https://backup.ionos.com (or via DCD, https://dcd.ionos.com/latest/). Use `backupunit get-sso-url` for a one-click SSO link.

Required values to run a command:

* Name
* Email
* Password

## Options

```text
  -u, --api-url string    Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [BackupUnitId Name Email State]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -e, --email string      E-mail address that will receive backup service reports. Does not need to match your IONOS CLOUD API username (required)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -n, --name string       Alphanumeric name for the BackupUnit. Combined with your contract number it forms the backup login (CONTRACT_NUMBER-NAME), so it must be globally unique and CANNOT be changed after creation (required)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -p, --password string   Login secret used to register the backup agent. Write-only: it is never returned by the API, so record it now (required)
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a BackupUnit
ionosctl compute backupunit create --name mybackups --email ops@example.com --password 'S3cretPass!'

# Then open the backup console via SSO (grab the id from the create output)
ionosctl compute backupunit get-sso-url --backupunit-id BACKUPUNIT_ID
```

