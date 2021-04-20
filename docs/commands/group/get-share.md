---
description: Get a Resource Share from a Group
---

# GetShare

## Usage

```text
ionosctl group get-share [flags]
```

## Description

Use this command to retrieve the details of a specific Shared Resource available to a specified Group.

.

Required values to run command:

* Group Id
* Resource Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Columns to be printed in the standard output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --group-id string      The unique Group Id [Required flag]
  -h, --help                 help for get-share
      --ignore-stdin         Force command to execute without user input
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
      --resource-id string   The unique Resource Id [Required flag]
```

## Examples

```text
ionosctl group get-share --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f --resource-id cefc2175-001f-4b94-8693-6263d731fe8e 
ShareId                                EditPrivilege   SharePrivilege
cefc2175-001f-4b94-8693-6263d731fe8e   false           true
```

