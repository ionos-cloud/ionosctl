---
description: List Users from a Group
---

# GroupUserList

## Usage

```text
ionosctl group user list [flags]
```

## Description

Use this command to get a list of Users from a specific Group.

Required values to run command:

* Group Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Columns to be printed in the standard output (default [GroupId,Name,CreateDataCenter,CreateSnapshot,ReserveIp,AccessActivityLog,CreatePcc,S3Privilege,CreateBackupUnit,CreateInternetAccess,CreateK8s])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
      --group-id string   The unique Group Id (required)
  -h, --help              help for list
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
```

## Examples

```text
ionosctl group user list --group-id 45ba215b-6897-40b6-879c-cbadb527cefd 
UserId                                 Firstname   Lastname   Email                    S3CanonicalUserId                  Administrator   ForceSecAuth   SecAuthActive   Active
62599641-aa2d-4ecc-bdc4-118f5f39f23d   test        test       testrandom53@gmail.com   f670112b3e74038b51db78d5836d7854   false           false          false           true
```

