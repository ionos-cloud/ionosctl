---
description: Get a Group
---

# GroupGet

## Usage

```text
ionosctl group get [flags]
```

## Aliases

For `group` command:

```text
[g]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific Group.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [GroupId Name CreateDataCenter CreateSnapshot ReserveIp AccessActivityLog CreatePcc S3Privilege CreateBackupUnit CreateInternetAccess CreateK8s CreateFlowLog AccessAndManageMonitoring AccessAndManageCertificates] (default [GroupId,Name,CreateDataCenter,CreateSnapshot,CreatePcc,CreateBackupUnit,CreateInternetAccess,CreateK8s,ReserveIp])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -i, --group-id string   The unique Group Id (required)
  -h, --help              help for get
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           see step by step process when running a command
```

## Examples

```text
ionosctl group get --group-id GROUP_ID
```

