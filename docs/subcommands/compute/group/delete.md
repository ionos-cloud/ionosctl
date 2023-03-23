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
  -a, --all                Delete all Groups.
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [GroupId Name CreateDataCenter CreateSnapshot ReserveIp AccessActivityLog CreatePcc S3Privilege CreateBackupUnit CreateInternetAccess CreateK8s CreateFlowLog AccessAndManageMonitoring AccessAndManageCertificates] (default [GroupId,Name,CreateDataCenter,CreateSnapshot,CreatePcc,CreateBackupUnit,CreateInternetAccess,CreateK8s,ReserveIp])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -f, --force              Force command to execute without user input
  -i, --group-id string    The unique Group Id (required)
  -h, --help               Print usage
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
  -t, --timeout int        Timeout option for Request for Group deletion [seconds] (default 60)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait-for-request   Wait for Request for Group deletion to be executed
```

## Examples

```text
ionosctl group delete --group-id GROUP_ID
```

