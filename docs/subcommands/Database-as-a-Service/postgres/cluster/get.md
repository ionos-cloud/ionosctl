---
description: "Get a PostgreSQL Cluster"
---

# DbaasPostgresClusterGet

## Usage

```text
ionosctl dbaas postgres cluster get [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
```

For `cluster` command:

```text
[c]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a PostgreSQL Cluster by using its ID.

Required values to run command:

* Cluster Id

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [ClusterId DisplayName Location State PostgresVersion Instances Ram Cores StorageSize StorageType DatacenterId LanId Cidr MaintenanceWindow SynchronizationMode BackupLocation] (default [ClusterId,DisplayName,Location,DatacenterId,LanId,Cidr,Instances,State])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout option for Cluster to be in AVAILABLE state [seconds] (default 1200)
  -v, --verbose             Print step-by-step process when running command
  -W, --wait-for-state      Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl dbaas postgres cluster get -i CLUSTER_ID
```

