---
description: "Delete a Group"
---

# GroupDelete

## Usage

```text
ionosctl compute group delete [flags]
```

## Aliases

For `group` command:

```text
[g]
```

For `delete` command:

```text
[d]
```

## Description

Delete a single Group. This removes the Group and the privileges it granted; its members lose those privileges (unless another Group still grants them) and lose access to the Group's shared resources.

Deleting a Group does NOT delete the Users in it, nor the resources shared with it - the resources simply become inaccessible to former members, except to a Contract Owner, Admin (administrator user), or Resource Owner.

Required values to run command:

* Group Id

## Options

```text
  -a, --all               Delete every Group on the contract. Use with caution: this strips group-granted privileges and shared-resource access from all members
  -u, --api-url string    Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [GroupId Name CreateDataCenter CreateSnapshot CreatePcc CreateBackupUnit CreateInternetAccess CreateK8s ReserveIp AccessActivityLog S3Privilege CreateFlowLog AccessAndManageMonitoring AccessAndManageCertificates AccessAndManageDns ManageDBaaS ManageRegistry]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -i, --group-id string   The unique Group Id (required)
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute group delete --group-id GROUP_ID
```

