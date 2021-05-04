---
description: List Groups
---

# List

## Usage

```text
ionosctl group list [flags]
```

## Description

Use this command to get a list of available Groups available on your account.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl group list
GroupId                                Name   CreateDataCenter   CreateSnapshot   ReserveIp   AccessActivityLog   CreatePcc   S3Privilege   CreateBackupUnit   CreateInternetAccess   CreateK8s
1d500d7a-43af-488a-a656-79e902433767   test   false              false            false       false               false       false         false              false                  false
```

