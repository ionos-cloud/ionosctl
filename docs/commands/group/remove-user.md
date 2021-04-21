---
description: Remove User from a Group
---

# RemoveUser

## Usage

```text
ionosctl group remove-user [flags]
```

## Description

Use this command to remove a User from a Group.

Required values to run command:

* Group Id
* User Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Columns to be printed in the standard output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --group-id string   The unique Group Id [Required flag]
  -h, --help              help for remove-user
      --ignore-stdin      Force command to execute without user input
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
      --user-id string    The unique User Id [Required flag]
```

## Examples

```text
ionosctl group remove-user --group-id e99f4cdb-746d-4c3c-b38c-b749ca23f917 --user-id e3b5fefb-27b0-4eea-a7c4-c57934ad23cb --wait 
Warning: Are you sure you want to remove user from group (y/N) ? 
y
Waiting for request: 353eff98-120f-4d91-82f5-a8aff1ddb277
RequestId: 353eff98-120f-4d91-82f5-a8aff1ddb277
Status: Command group remove-user and request have been successfully executed
```

