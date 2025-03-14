---
description: "Get a Group"
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
                          Available columns: [GroupId Name CreateDataCenter CreateSnapshot ReserveIp AccessActivityLog CreatePcc S3Privilege CreateBackupUnit CreateInternetAccess CreateK8s CreateFlowLog AccessAndManageMonitoring AccessAndManageCertificates AccessAndManageDns ManageDBaaS ManageRegistry ManageDataplatform] (default [GroupId,Name,CreateDataCenter,CreateSnapshot,CreatePcc,CreateBackupUnit,CreateInternetAccess,CreateK8s,ReserveIp])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32       Controls the detail depth of the response objects. Max depth is 10.
  -f, --force             Force command to execute without user input
  -i, --group-id string   The unique Group Id (required)
  -h, --help              Print usage
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for polling the request (default 60)
  -v, --verbose           Print step-by-step process when running command
  -w, --wait              Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl group get --group-id GROUP_ID
```

