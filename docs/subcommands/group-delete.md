---
description: Delete a Group
---

# GroupDelete

## Usage

```text
ionosctl group delete [flags]
```

## Description

Use this operation to delete a single Group. Resources that are assigned to the Group are NOT deleted, but are no longer accessible to the Group members unless the member is a Contract Owner, Admin, or Resource Owner.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string     Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -F, --format strings     Collection of fields to be printed on output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
      --group-id string    The unique Group Id (required)
  -h, --help               help for delete
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
  -t, --timeout int        Timeout option for Request for Group deletion [seconds] (default 60)
  -w, --wait-for-request   Wait for Request for Group deletion to be executed
```

## Examples

```text
ionosctl group delete --group-id 1d500d7a-43af-488a-a656-79e902433767 
Warning: Are you sure you want to delete group (y/N) ? 
y
RequestId: e20d2851-0d20-453d-b752-ed1c34a83625
Status: Command group delete has been successfully executed
```

