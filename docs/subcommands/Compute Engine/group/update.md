---
description: "Update a Group"
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
      --access-certs        Privilege for a group to access and manage certificates. E.g.: --access-certs=true, --access-certs=false
      --access-dns          Privilege for a group to access and manage dns records
      --access-logs         The group will be allowed to access the activity log. E.g.: --access-logs=true, --access-logs=false
      --access-monitoring   Privilege for a group to access and manage monitoring related functionality using Monotoring-as-a-Service. E.g.: --access-monitoring=true, --access-monitoring=false
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [GroupId Name CreateDataCenter CreateSnapshot ReserveIp AccessActivityLog CreatePcc S3Privilege CreateBackupUnit CreateInternetAccess CreateK8s CreateFlowLog AccessAndManageMonitoring AccessAndManageCertificates AccessAndManageDns ManageDBaaS ManageRegistry] (default [GroupId,Name,CreateDataCenter,CreateSnapshot,CreatePcc,CreateBackupUnit,CreateInternetAccess,CreateK8s,ReserveIp])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --create-backup       The group will be able to manage Backup Units. E.g.: --create-backup=true, --create-backup=false
      --create-dc           The group will be allowed to create Data Centers. E.g.: --create-dc=true, --create-dc=false
      --create-flowlog      The group will be allowed to create Flow Logs. E.g.: --create-flowlog=true, --create-flowlog=false
      --create-k8s          The group will be allowed to create K8s Clusters. E.g.: --create-k8s=true, --create-k8s=false
      --create-nic          The group will be allowed to create NICs. E.g.: --create-nic=true, --create-nic=false
      --create-pcc          The group will be allowed to create PCCs. E.g.: --create-pcc=true, --create-pcc=false
      --create-snapshot     The group will be allowed to create Snapshots. E.g.: --create-snapshot=true, --create-snapshot=false
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10.
  -f, --force               Force command to execute without user input
  -i, --group-id string     The unique Group Id (required)
  -h, --help                Print usage
      --manage-dbaas        Privilege for a group to manage DBaaS related functionality
      --manage-registry     Privilege for group accessing container registry related functionality
  -n, --name string         Name for the Group
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --reserve-ip          The group will be allowed to reserve IP addresses. E.g.: --reserve-ip=true, --reserve-ip=false
      --s3privilege         The group will be allowed to manage S3. E.g.: --s3privilege=true, --s3privilege=false
  -t, --timeout int         Timeout option for Request for Group update [seconds] (default 60)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request    Wait for Request for Group update to be executed
```

## Examples

```text
ionosctl group update --group-id GROUP_ID --reserve-ip
```

