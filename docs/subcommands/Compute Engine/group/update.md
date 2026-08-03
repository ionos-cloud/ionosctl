---
description: "Update a Group's name or privileges"
---

# GroupUpdate

## Usage

```text
ionosctl compute group update [flags]
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

Update an existing IAM Group's name and/or privileges. Any privilege flag you set is applied to the Group and therefore to ALL of its members; privilege flags you do NOT pass keep their current value, so you can toggle a single capability without listing the rest.

Setting a flag to false REVOKES that capability from every member of the Group (unless another Group they belong to still grants it).

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Group Id

## Options

```text
      --access-certs        Grant the 'access and manage certificates' privilege: members may manage certificates in the Certificate Manager. E.g.: --access-certs=true, --access-certs=false
      --access-dns          Grant the 'access and manage DNS' privilege: members may manage DNS zones and records
      --access-logs         Grant the 'access activity log' privilege: members may read the contract's audit/activity log. E.g.: --access-logs=true, --access-logs=false
      --access-monitoring   Grant the 'access and manage monitoring' privilege: members may access metrics and manage alarms via Monitoring-as-a-Service (MaaS). E.g.: --access-monitoring=true, --access-monitoring=false
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [GroupId Name CreateDataCenter CreateSnapshot CreatePcc CreateBackupUnit CreateInternetAccess CreateK8s ReserveIp AccessActivityLog S3Privilege CreateFlowLog AccessAndManageMonitoring AccessAndManageCertificates AccessAndManageDns ManageDBaaS ManageRegistry]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --create-backup       Grant the 'create backup unit' privilege: members may create and manage Backup Units. E.g.: --create-backup=true, --create-backup=false
      --create-dc           Grant the 'create data center' privilege: members may create new Virtual Data Centers. E.g.: --create-dc=true, --create-dc=false
      --create-flowlog      Grant the 'create Flow Log' privilege: members may create Flow Logs to capture network traffic. E.g.: --create-flowlog=true, --create-flowlog=false
      --create-k8s          Grant the 'create Kubernetes cluster' privilege: members may create Managed Kubernetes clusters. E.g.: --create-k8s=true, --create-k8s=false
      --create-nic          Grant the 'create internet access' privilege: members may attach public/internet-facing connectivity (despite the flag name, this is NOT a per-NIC toggle). E.g.: --create-nic=true, --create-nic=false
      --create-pcc          Grant the 'create PCC' privilege: members may create Private Cross-Connects to bridge private LANs across datacenters. E.g.: --create-pcc=true, --create-pcc=false
      --create-snapshot     Grant the 'create snapshot' privilege: members may take snapshots of volumes. E.g.: --create-snapshot=true, --create-snapshot=false
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -i, --group-id string     The unique Group Id (required)
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
      --manage-dbaas        Grant the 'manage DBaaS' privilege: members may manage Database-as-a-Service clusters (PostgreSQL, MongoDB, MariaDB, etc.)
      --manage-registry     Grant the 'manage Registry' privilege: members may manage Container Registry repositories
  -n, --name string         The new name for the Group
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
      --reserve-ip          Grant the 'reserve IP' privilege: members may reserve public IP blocks. E.g.: --reserve-ip=true, --reserve-ip=false
      --s3privilege         Grant the S3 privilege: members may use IONOS Object Storage (S3-compatible) and manage their own S3 keys. E.g.: --s3privilege=true, --s3privilege=false
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a group
ionosctl compute group update --group-id GROUP_ID --name "Platform Team"

# Grant an extra capability without touching the group's other privileges
ionosctl compute group update --group-id GROUP_ID --reserve-ip

# Revoke the ability to create datacenters and Kubernetes clusters
ionosctl compute group update --group-id GROUP_ID --create-dc=false --create-k8s=false
```

