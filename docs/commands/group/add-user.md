---
description: Add User to a Group
---

# AddUser

## Usage

```text
ionosctl group add-user [flags]
```

## Description

Use this command to add an existing User to a specific Group.

Required values to run command:

* Group Id
* User Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Columns to be printed in the standard output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force             Force command to execute without user input
      --group-id string   The unique Group Id [Required flag]
  -h, --help              help for add-user
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
      --user-id string    The unique User Id [Required flag]
```

## Examples

```text
ionosctl group add-user --group-id e99f4cdb-746d-4c3c-b38c-b749ca23f917 --user-id 53d68de9-931a-4b61-b532-82f7b27afef3
UserId                                 Firstname   Lastname   Email                    Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId                  Active
53d68de9-931a-4b61-b532-82f7b27afef3   test1       test1      testrandom13@ionos.com   false           false          false           8b9dd6f39e613adb7a837127edb67d38   true
```

