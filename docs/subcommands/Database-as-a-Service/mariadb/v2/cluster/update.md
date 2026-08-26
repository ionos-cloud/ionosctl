---
description: "Update a MariaDB Cluster"
---

# DbaasMariadbV2ClusterUpdate

## Usage

```text
ionosctl dbaas mariadb-v2 cluster update [flags]
```

## Aliases

For `mariadb-v2` command:

```text
[maria-v2 mar-v2 ma-v2]
```

For `cluster` command:

```text
[c]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update attributes of a MariaDB Cluster. This command uses a combination of GET and PUT to simulate a PATCH operation.

Note the API's sizing constraints: instances and storage size can only be increased (never decreased), the version can only be upgraded (no downgrade), while cores and RAM can be both increased and decreased.

Required values to run command:

* Cluster Id

## Options

```text
  -u, --api-url string                Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mariadbv2' and env var 'IONOS_API_URL' (default "https://mariadb.%s.ionos.com")
      --backup-retention-days int32   Configures how many days cluster backups are retained. Minimum: 1, Maximum: 365
  -i, --cluster-id string             The unique ID of the Cluster (required)
      --cols strings                  Set of columns to be printed on output 
                                      Available columns: [ClusterId Name DnsName Version Instances State Cores Ram StorageSize Description BackupLocation RetentionDays MaintenanceDay MaintenanceTime LogsEnabled MetricsEnabled DatacenterId LanId Cidr StatusMessage]
  -c, --config string                 Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int32                   The number of CPU cores per instance. Can be increased or decreased
      --database string               Database for the credentials. Defaults to the cluster's current database
  -D, --depth int                     Level of detail for response objects (default 1)
  -F, --filters strings               Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                         Force command to execute without user input
  -h, --help                          Print usage
      --instances int32               The total number of instances in the cluster (one primary and n-1 secondary). Can only be increased
      --limit int                     Maximum number of items to return per request (default 50)
  -l, --location string               Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --logs-enabled                  Enable collection and reporting of logs for this cluster
      --maintenance-day string        Day of the week for the MaintenanceWindow. Must be specified together with --maintenance-time. Can be one of: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday
      --maintenance-time string       Time for the MaintenanceWindow. e.g.: 16:30:59. Must be specified together with --maintenance-day
      --metrics-enabled               Enable collection and reporting of metrics for this cluster
  -n, --name string                   The friendly name of your cluster
      --no-headers                    Don't print table headers when table output is used
      --offset int                    Number of items to skip before starting to collect the results
      --order-by string               Property to order the results by
  -o, --output string                 Desired output format [text|json|api-json] (default "text")
      --password string               Password for the database user. Required because the API does not return it on GET requests (minimum length 10) (required)
      --query string                  JMESPath query string to filter the output
  -q, --quiet                         Quiet output
      --ram string                    The amount of memory per instance. e.g. --ram 4, --ram 4GB. Can be increased or decreased
      --storage-size string           The size of the storage per instance. e.g. --storage-size 10GB. Can only be increased
  -t, --timeout int                   Timeout in seconds for --wait and other wait operations (default 600)
      --user string                   Username for the database user. Defaults to the cluster's current username
  -v, --verbose count                 Increase verbosity level [-v, -vv, -vvv]
      --version string                The MariaDB version of your cluster. Downgrades are not supported
  -w, --wait                          Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mariadb-v2 cluster update --location <location> --cluster-id <cluster-id> --password <password> --cores 4 --ram 8GB
```

