---
description: Create a BackupUnit
---

# BackupunitCreate

## Usage

```text
ionosctl backupunit create [flags]
```

## Description

Use this command to create a BackupUnit under a particular contract. You need to specify the name, email and password for the new BackupUnit.

Notes:

* The name assigned to the BackupUnit will be concatenated with the contract number to create the login name for the backup system. The name may NOT be changed after creation.
* The password set here is used along with the login name described above to register the backup agent with the backup system. When setting the password, please make a note of it, as the value cannot be retrieved using the Cloud API.
* The e-mail address supplied here does NOT have to be the same as your Cloud API username. This e-mail address will receive service reports from the backup system.
* To login to backup agent, please use https://dcd.ionos.com/latest/ and access BackupUnit Console or use https://backup.ionos.com

Required values to run a command:

* BackupUnit Name
* BackupUnit Email
* BackupUnit Password

## Options

```text
  -u, --api-url string               Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --backupunit-email string      The e-mail address you want to assign to the BackupUnit (required)
      --backupunit-name string       Alphanumeric name you want to assign to the BackupUnit (required)
      --backupunit-password string   Alphanumeric password you want to assign to the BackupUnit (required)
      --cols strings                 Columns to be printed in the standard output (default [BackupUnitId,Name,Email,State])
  -c, --config string                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                        Force command to execute without user input
  -h, --help                         help for create
  -o, --output string                Desired output format [text|json] (default "text")
  -q, --quiet                        Quiet output
      --timeout int                  Timeout option for Request for BackupUnit creation [seconds] (default 60)
      --wait-for-request             Wait for the Request for BackupUnit creation to be executed
```

## Examples

```text
ionosctl backupunit create --backupunit-name test1234test --backupunit-email testrandom18@ionos.com --backupunit-password ********
NOTE: To login with backup agent use: https://backup.ionos.com, with CONTRACT_NUMBER-BACKUP_UNIT_NAME and BACKUP_UNIT_PASSWORD!
BackupUnitId                           Name           Email
271a0627-70eb-4e36-8ff5-2e190f88cd2b   test1234test   testrandom18@ionos.com
RequestId: 2cd34841-f0b1-4ac7-9741-89a2575a9962
Status: Command backupunit create has been successfully executed
```

