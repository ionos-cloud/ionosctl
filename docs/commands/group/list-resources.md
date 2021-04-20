---
description: List Resources from a Group
---

# ListResources

## Usage

```text
ionosctl group list-resources [flags]
```

## Description

Use this command to get a list of Resources assigned to a Group. To see more details about Resources under a specific User, use `ionosctl user` commands.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Columns to be printed in the standard output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --group-id string   The unique Group Id [Required flag]
  -h, --help              help for list-resources
      --ignore-stdin      Force command to execute without user input
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
```

## Examples

```text
ionosctl group list-resources --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f 
ResourceId                             Name   SecAuthProtection   Type
cefc2175-001f-4b94-8693-6263d731fe8e          false               datacenter
```

