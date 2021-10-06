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

For `create` command:

```text
[c]
```

## Description

Use this command to create a new Group and set Group privileges. You can specify the name for the new Group. By default, all privileges will be set to false. You need to use flags privileges to be set to true.

## Options

```text
      --access-certs        Privilege for a group to access and manage certificates
      --access-logs         The group will be allowed to access the activity log
      --access-monitoring   Privilege for a group to access and manage monitoring related functionality using Monotoring-as-a-Service
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [GroupId Name CreateDataCenter CreateSnapshot ReserveIp AccessActivityLog CreatePcc S3Privilege CreateBackupUnit CreateInternetAccess CreateK8s CreateFlowLog AccessAndManageMonitoring AccessAndManageCertificates] (default [GroupId,Name,CreateDataCenter,CreateSnapshot,CreatePcc,CreateBackupUnit,CreateInternetAccess,CreateK8s,ReserveIp])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --create-backup       The group will be able to manage Backup Units
      --create-dc           The group will be allowed to create Data Centers
      --create-flowlog      The group will be allowed to create Flow Logs
      --create-k8s          The group will be allowed to create K8s Clusters
      --create-nic          The group will be allowed to create NICs
      --create-pcc          The group will be allowed to create PCCs
      --create-snapshot     The group will be allowed to create Snapshots
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -n, --name string         Name for the Group (default "Unnamed Group")
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
      --reserve-ip          The group will be allowed to reserve IP addresses
      --s3privilege         The group will be allowed to manage S3
  -t, --timeout int         Timeout option for Request for Group creation [seconds] (default 60)
  -v, --verbose             Print step-by-step process when running command
  -w, --wait-for-request    Wait for Request for Group creation to be executed
```

## Examples

```text
ionosctl group create --name NAME --wait-for-request
```

