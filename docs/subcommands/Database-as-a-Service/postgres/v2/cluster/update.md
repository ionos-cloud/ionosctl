---
description: "Update a PostgreSQL Cluster"
---

# DbaasPostgresV2ClusterUpdate

## Usage

```text
ionosctl dbaas postgres-v2 cluster update [flags]
```

## Aliases

For `postgres-v2` command:

```text
[pg-v2 pgsql-v2 postgresql-v2 psql-v2]
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

Use this command to update attributes of a PostgreSQL Cluster.

Required values to run command:

* Cluster Id
* DB Password

## Options

```text
  -u, --api-url string             Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'psqlv2' and env var 'IONOS_API_URL' (default "https://postgresql.%s.ionos.com")
  -C, --cidr string                The IP and subnet for the cluster. Note the following unavailable IP range: 10.208.0.0/12. e.g.: 192.168.1.100/24
  -i, --cluster-id string          The unique ID of the Cluster (required)
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [ClusterId DisplayName DnsName PostgresVersion Instances Ram Cores StorageSize State SyncMode Description ConnectionPooler MaintenanceDay MaintenanceTime BackupLocation LogsEnabled MetricsEnabled DatacenterId LanId Cidr DbUsername DbDatabase StatusMessage]
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --connection-pooler string   Connection pooling mode. Can be one of: DISABLED, TRANSACTION, SESSION
      --cores int                  The number of CPU cores per instance. Minimum: 1, Maximum: 62
      --datacenter-id string       The unique ID of the Datacenter to connect to your cluster. It has to be in the same location as the current datacenter
  -P, --db-password string         Password for the initial postgres user. Required because the API does not return it on GET requests (required)
  -D, --depth int                  Level of detail for response objects (default 1)
      --description string         Human-readable description for the cluster
  -F, --filters strings            Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
  -I, --instances int              The number of instances in your cluster. Minimum: 1, Maximum: 5
  -L, --lan-id string              The unique ID of the LAN to connect your cluster to
      --limit int                  Maximum number of items to return per request (default 50)
  -l, --location string            Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, gb/bhx, us/las, us/mci, us/ewr (default "de/txl")
      --logs-enabled               Enable collection and reporting of logs for this cluster
  -d, --maintenance-day string     Day of the week for the MaintenanceWindow. Must be specified together with --maintenance-time. Can be one of: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday
  -T, --maintenance-time string    Time for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59. Must be specified together with --maintenance-day
      --metrics-enabled            Enable collection and reporting of metrics for this cluster
  -n, --name string                The friendly name of your cluster
      --no-headers                 Don't print table headers when table output is used
      --offset int                 Number of items to skip before starting to collect the results
      --order-by string            Property to order the results by
  -o, --output string              Desired output format [text|json|api-json] (default "text")
      --query string               JMESPath query string to filter the output
  -q, --quiet                      Quiet output
      --ram string                 The amount of memory per instance in GB. Minimum: 4, Maximum: 240. e.g. --ram 4, --ram 4GB
      --storage-size string        The amount of storage per instance in GB. Minimum: 10, Maximum: 4096. e.g.: --storage-size 20, --storage-size 20GB
  -S, --sync-mode string           Replication mode. Can be one of: ASYNCHRONOUS, STRICTLY_SYNCHRONOUS
  -t, --timeout int                Timeout option for Cluster to be in AVAILABLE state[seconds] (default 1200)
  -v, --verbose count              Increase verbosity level [-v, -vv, -vvv]
  -V, --version string             The PostgreSQL version of your cluster
  -W, --wait-for-state             Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl dbaas postgres-v2 cluster update --cluster-id <cluster-id> --db-password <password> --cores 4 --ram 8GB
```

