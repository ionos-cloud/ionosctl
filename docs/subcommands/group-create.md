---
description: Create a Group
---

# GroupCreate

## Usage

```text
ionosctl group create [flags]
```

## Aliases

For `group` command:
```text
[g]
```

## Description

Use this command to create a new Group and set Group privileges. You need to specify the name for the new Group. By default, all privileges will be set to false. You need to use flags privileges to be set to true.

Required values to run a command:

* Name

## Options

```text
      --access-logs        The group will be allowed to access the activity log
  -u, --api-url string     Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --create-backup      The group will be able to manage Backup Units
      --create-dc          The group will be allowed to create Data Centers
      --create-k8s         The group will be allowed to create K8s Clusters
      --create-nic         The group will be allowed to create NICs
      --create-pcc         The group will be allowed to create PCCs
      --create-snapshot    The group will be allowed to create Snapshots
  -f, --force              Force command to execute without user input
  -F, --format strings     Collection of fields to be printed on output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -h, --help               help for create
  -n, --name string        Name for the Group (required)
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
      --reserve-ip         The group will be allowed to reserve IP addresses
      --s3privilege        The group will be allowed to manage S3
  -t, --timeout int        Timeout option for Request for Group creation [seconds] (default 60)
  -w, --wait-for-request   Wait for Request for Group creation to be executed
```

## Examples

```text
ionosctl group create --name test --wait-for-request
1.2s Waiting for request... DONE
GroupId                                Name   CreateDataCenter   CreateSnapshot   ReserveIp   AccessActivityLog   CreatePcc   S3Privilege   CreateBackupUnit   CreateInternetAccess   CreateK8s
1d500d7a-43af-488a-a656-79e902433767   test   false              false            false       false               false       false         false              false                  false
```

