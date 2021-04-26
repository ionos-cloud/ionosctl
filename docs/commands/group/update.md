---
description: Update a Group
---

# Update

## Usage

```text
ionosctl group update [flags]
```

## Description

Use this command to update details about a specific Group.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string          Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings            Columns to be printed in the standard output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string           Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                   Force command to execute without user input
      --group-access-logs       The group will be allowed to access the activity log
      --group-create-backup     The group will be able to manage Backup Units
      --group-create-dc         The group will be allowed to create Data Centers
      --group-create-k8s        The group will be allowed to create K8s Clusters
      --group-create-nic        The group will be allowed to create NICs
      --group-create-pcc        The group will be allowed to create PCCs
      --group-create-snapshot   The group will be allowed to create Snapshots
      --group-id string         The unique Group Id [Required flag]
      --group-name string       Name for the Group [Required flag]
      --group-reserve-ip        The group will be allowed to reserve IP addresses
      --group-s3privilege       The group will be allowed to manage S3
  -h, --help                    help for update
  -o, --output string           Desired output format [text|json] (default "text")
  -q, --quiet                   Quiet output
      --timeout int             Timeout option for Group to be updated [seconds] (default 60)
      --wait                    Wait for Group attributes to be updated
```

## Examples

```text
ionosctl group update --group-id e99f4cdb-746d-4c3c-b38c-b749ca23f917 --group-reserve-ip 
GroupId                                Name         CreateDataCenter   CreateSnapshot   ReserveIp   AccessActivityLog   CreatePcc   S3Privilege   CreateBackupUnit   CreateInternetAccess   CreateK8s
e99f4cdb-746d-4c3c-b38c-b749ca23f917   testUpdate   true               true             true        false               false       false         false              false                  true
RequestId: 2bfe43a4-ea09-48fc-bb53-136c7f7d061f
Status: Command group update has been successfully executed
```

