---
description: "Create a BackupUnit"
---

# BackupunitCreate

## Usage

```text
ionosctl backupunit create [flags]
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

Use this command to create a BackupUnit under a particular contract. You need to specify the name, email and password for the new BackupUnit.

Notes:

* The name assigned to the BackupUnit will be concatenated with the contract number to create the login name for the backup system. The name may NOT be changed after creation.
* The password set here is used along with the login name described above to register the backup agent with the backup system. When setting the password, please make a note of it, as the value cannot be retrieved using the Cloud API.
* The e-mail address supplied here does NOT have to be the same as your Cloud API username. This e-mail address will receive service reports from the backup system.
* To login to backup agent, please use [https://dcd.ionos.com/latest/](https://dcd.ionos.com/latest/) and access BackupUnit Console or use [https://backup.ionos.com](https://backup.ionos.com)

Required values to run a command:

* Name
* Email
* Password

## Options

```text
      --cols strings      Set of columns to be printed on output 
                          Available columns: [BackupUnitId Name Email State] (default [BackupUnitId,Name,Email,State])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32       Controls the detail depth of the response objects. Max depth is 10.
  -e, --email string      The e-mail address you want to assign to the BackupUnit (required)
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -n, --name string       Alphanumeric name you want to assign to the BackupUnit (required)
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -p, --password string   Alphanumeric password you want to assign to the BackupUnit (required)
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout option for Request for BackupUnit creation [seconds] (default 60)
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl backupunit create --name NAME --email EMAIL --password PASSWORD
```

