---
description: List Groups
---

# GroupList

## Usage

```text
ionosctl group list [flags]
```

## Aliases

For `group` command:

```text
[g]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of available Groups available on your account

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name createDataCenter createSnapshot reserveIp accessActivityLog createPcc s3Privilege createBackupUnit createInternetAccess createK8sCluster createFlowLog accessAndManageMonitoring accessAndManageCertificates manageDBaaS]

## Options

```text
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings     Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          When using text output, don't print headers
      --order-by string     Limits results to those containing a matching value for a specific property
```

## Examples

```text
ionosctl group list
```

