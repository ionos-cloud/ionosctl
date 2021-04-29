---
description: Create a Group
---

# Create

## Usage

```text
ionosctl group create [flags]
```

## Description

Use this command to create a new Group and set Group privileges. You need to specify the name for the new Group. By default, all privileges will be set to false. You need to use flags privileges to be set to true.

Required values to run a command:

* Group Name

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
      --group-name string       Name for the Group (required)
      --group-reserve-ip        The group will be allowed to reserve IP addresses
      --group-s3privilege       The group will be allowed to manage S3
  -h, --help                    help for create
  -o, --output string           Desired output format [text|json] (default "text")
  -q, --quiet                   Quiet output
      --timeout int             Timeout option for Group to be created [seconds] (default 60)
      --wait                    Wait for Group to be created
```

## Examples

```text
ionosctl group create --group-name test --wait 
Waiting for request: eae6bb8b-3736-4cf0-bc71-72a95d1b2a63
GroupId                                Name   CreateDataCenter   CreateSnapshot   ReserveIp   AccessActivityLog   CreatePcc   S3Privilege   CreateBackupUnit   CreateInternetAccess   CreateK8s
1d500d7a-43af-488a-a656-79e902433767   test   false              false            false       false               false       false         false              false                  false
```

