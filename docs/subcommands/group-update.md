---
description: Update a Group
---

# GroupUpdate

## Usage

```text
ionosctl group update [flags]
```

## Aliases

For `group` command:

```text
[g]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update details about a specific Group.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Group Id

## Options

```text
      --access-logs        The group will be allowed to access the activity log
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [GroupId Name CreateDataCenter CreateSnapshot ReserveIp AccessActivityLog CreatePcc S3Privilege CreateBackupUnit CreateInternetAccess CreateK8s] (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --create-backup      The group will be able to manage Backup Units
      --create-dc          The group will be allowed to create Data Centers
      --create-k8s         The group will be allowed to create K8s Clusters
      --create-nic         The group will be allowed to create NICs
      --create-pcc         The group will be allowed to create PCCs
      --create-snapshot    The group will be allowed to create Snapshots
  -f, --force              Force command to execute without user input
  -i, --group-id string    The unique Group Id (required)
  -h, --help               Print usage
  -n, --name string        Name for the Group
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
      --reserve-ip         The group will be allowed to reserve IP addresses
      --s3privilege        The group will be allowed to manage S3
  -t, --timeout int        Timeout option for Request for Group update [seconds] (default 60)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait-for-request   Wait for Request for Group update to be executed
```

## Examples

```text
ionosctl group update --group-id GROUP_ID --reserve-ip
```

