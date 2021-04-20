---
description: Get a Group
---

# Get

## Usage

```text
ionosctl group get [flags]
```

## Description

Use this command to retrieve details about a specific Group.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Columns to be printed in the standard output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --group-id string   The unique Group Id [Required flag]
  -h, --help              help for get
      --ignore-stdin      Force command to execute without user input
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
```

## Examples

```text
ionosctl group get --group-id 1d500d7a-43af-488a-a656-79e902433767 
GroupId                                Name   CreateDataCenter   CreateSnapshot   ReserveIp   AccessActivityLog   CreatePcc   S3Privilege   CreateBackupUnit   CreateInternetAccess   CreateK8s
1d500d7a-43af-488a-a656-79e902433767   test   false              false            false       false               false       false         false              false                  false
```

