---
description: Remove a Resource Share from a Group
---

# RemoveShare

## Usage

```text
ionosctl group remove-share [flags]
```

## Description

This command removes a Resource Share from a specified Group.

Required values to run command:

* Resource Id
* Group Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Columns to be printed in the standard output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --group-id string      The unique Group Id [Required flag]
  -h, --help                 help for remove-share
      --ignore-stdin         Force command to execute without user input
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
      --resource-id string   The unique Resource Id [Required flag]
      --timeout int          Timeout option for a Resource Share to be remove from a Group [seconds] (default 60)
      --wait                 Wait for a Resource Share to be remove from a Group
```

## Examples

```text
ionosctl group remove-share --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f --resource-id cefc2175-001f-4b94-8693-6263d731fe8e --wait 
Warning: Are you sure you want to remove share from group (y/N) ? 
y
Waiting for request: 9ff7e57f-b568-4257-b27f-13a4cf11a7fc
RequestId: 9ff7e57f-b568-4257-b27f-13a4cf11a7fc
Status: Command group remove-share and request have been successfully executed
```

