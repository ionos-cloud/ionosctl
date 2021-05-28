---
description: Delete a Group
---

# GroupDelete

## Usage

```text
ionosctl group delete [flags]
```

## Aliases

For `group` command:
```text
[g]
```

For `delete` command:
```text
[d]
```

## Description

Use this operation to delete a single Group. Resources that are assigned to the Group are NOT deleted, but are no longer accessible to the Group members unless the member is a Contract Owner, Admin, or Resource Owner.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string     Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [GroupId Name CreateDataCenter CreateSnapshot ReserveIp AccessActivityLog CreatePcc S3Privilege CreateBackupUnit CreateInternetAccess CreateK8s] (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -i, --group-id string    The unique Group Id (required)
  -h, --help               help for delete
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
  -t, --timeout int        Timeout option for Request for Group deletion [seconds] (default 60)
  -w, --wait-for-request   Wait for Request for Group deletion to be executed
```

## Examples

```text
ionosctl group delete --group-id GROUP_ID
```

